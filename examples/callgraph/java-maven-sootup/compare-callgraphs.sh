#!/bin/bash
# ============================================================
# compare-callgraphs.sh
#
# Compares the Soot v1 (classic) and SootUp v2 callgraph outputs.
# Handles v1 output in either plain JSON or Debricked's
# base64+zip post-processed format.
#
# Usage:
#   bash compare-callgraphs.sh [--run-both]
# ============================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
V1_DIR="$SCRIPT_DIR/../java-maven-minimal"
V2_DIR="$SCRIPT_DIR"

V1_OUTPUT="$V1_DIR/debricked-call-graph.java"
V2_OUTPUT="$V2_DIR/debricked-call-graph-sootup.java"

if [ "${1:-}" = "--run-both" ]; then
    echo "Running Soot v1..."
    (cd "$V1_DIR" && bash run-callgraph.sh)
    echo ""
    echo "Running SootUp v2..."
    (cd "$V2_DIR" && bash run-callgraph-sootup.sh)
    echo ""
fi

if [ ! -f "$V1_OUTPUT" ]; then
    echo "ERROR: Soot v1 output not found: $V1_OUTPUT"
    echo "  Run: cd $V1_DIR && bash run-callgraph.sh"
    exit 1
fi
if [ ! -f "$V2_OUTPUT" ]; then
    echo "ERROR: SootUp v2 output not found: $V2_OUTPUT"
    echo "  Run: cd $V2_DIR && bash run-callgraph-sootup.sh"
    exit 1
fi

python3 - "$V1_OUTPUT" "$V2_OUTPUT" <<'PY'
import base64, io, json, sys, zipfile

v1_path, v2_path = sys.argv[1], sys.argv[2]


def load_callgraph(path):
    raw = open(path, "rb").read()
    try:
        data = json.loads(raw.decode("utf-8", errors="strict"))
        return data, "json"
    except Exception:
        pass
    text = raw.decode("utf-8", errors="ignore").strip()
    decoded = base64.b64decode(text)
    with zipfile.ZipFile(io.BytesIO(decoded), "r") as zf:
        payload = zf.read(zf.namelist()[0]).decode("utf-8")
        return json.loads(payload), "base64+zip"


def short_sig(sig):
    """
    Normalize to a short, package-free signature for cross-format comparison.
    'com.example.callgraph.App.main(String[])' -> 'App.main(String[])'
    'App.main(String[])'                        -> 'App.main(String[])'
    """
    paren = sig.index("(")
    before_paren = sig[:paren]
    parts = before_paren.split(".")
    # last two parts are ClassName.methodName (drop all package segments)
    short_method = ".".join(parts[-2:]) if len(parts) >= 2 else before_paren
    return short_method + sig[paren:]


def get_methods(data):
    return {e[0] for e in data.get("data", []) if isinstance(e, list) and e}


def get_user_methods(data):
    return {
        e[0] for e in data.get("data", [])
        if isinstance(e, list) and len(e) > 1 and e[1] is True
    }


def has_edge(data, caller_substr, callee_substr):
    """Return (found, line_number) searching by substring in both FQN and short form."""
    for entry in data.get("data", []):
        if not (isinstance(entry, list) and len(entry) >= 8):
            continue
        sig = entry[0]
        callers_list = entry[7]
        if callee_substr in sig or callee_substr in short_sig(sig):
            for c in callers_list:
                if isinstance(c, list) and len(c) >= 2:
                    caller_sig = c[0]
                    if caller_substr in caller_sig or caller_substr in short_sig(caller_sig):
                        return True, c[1]
    return False, -1


v1, v1_fmt = load_callgraph(v1_path)
v2, v2_fmt = load_callgraph(v2_path)

v1_all   = get_methods(v1)
v2_all   = get_methods(v2)
v1_user  = get_user_methods(v1)
v2_user  = get_user_methods(v2)

# Normalize to short sigs for cross-format intersection
v1_user_short = {short_sig(m) for m in v1_user}
v2_user_short = {short_sig(m) for m in v2_user}

common   = sorted(v1_user_short & v2_user_short)
only_v1  = sorted(v1_user_short - v2_user_short)
only_v2  = sorted(v2_user_short - v1_user_short)

print("=" * 64)
print(" Callgraph Comparison: Soot v1 vs SootUp v2")
print("=" * 64)
print()
print(f"  v1 output format  : {v1_fmt}")
print(f"  v2 output format  : {v2_fmt}")
print()
print(f"  {'Metric':<35}  {'Soot v1':>10}  {'SootUp v2':>10}")
print(f"  {'-'*35}  {'-'*10}  {'-'*10}")
print(f"  {'Total methods':<35}  {len(v1_all):>10}  {len(v2_all):>10}")
print(f"  {'User code methods':<35}  {len(v1_user):>10}  {len(v2_user):>10}")
print(f"  {'Library/JDK methods':<35}  {len(v1_all)-len(v1_user):>10}  {len(v2_all)-len(v2_user):>10}")
print()

print("-" * 64)
print(f" User methods in BOTH ({len(common)})")
print("-" * 64)
for m in common:
    print(f"  ✓  {m}")
print()

if only_v1:
    print("-" * 64)
    print(f" In Soot v1 ONLY ({len(only_v1)})  — missing from SootUp")
    print("-" * 64)
    for m in only_v1:
        print(f"  ✗  {m}")
    print()

if only_v2:
    print("-" * 64)
    print(f" In SootUp v2 ONLY ({len(only_v2)})  — extra vs Soot")
    print("-" * 64)
    for m in only_v2:
        print(f"  ★  {m}")
    print()

checks = [
    ("App.main → OrderService.placeOrder",           "App.main",            "OrderService.placeOrder"),
    ("App.main → LoggerUtil.log",                    "App.main",            "LoggerUtil.log"),
    ("OrderService.placeOrder → calculateTotal",     "OrderService.placeOrder", "PricingService.calculateTotal"),
    ("PricingService.calculateTotal → applyDiscount","calculateTotal",      "applyDiscount"),
    ("LoggerUtil.log → StringUtils.upperCase",       "LoggerUtil.log",      "StringUtils.upperCase"),
]

print("-" * 64)
print(" Key call edges")
print("-" * 64)
all_ok = True
for label, caller, callee in checks:
    ok1, ln1 = has_edge(v1, caller, callee)
    ok2, ln2 = has_edge(v2, caller, callee)
    s1 = f"✓ line {ln1}" if ok1 else "✗ MISSING"
    s2 = f"✓ line {ln2}" if ok2 else "✗ MISSING"
    match = "=" if ok1 == ok2 else "≠"
    if not ok1 or not ok2:
        all_ok = False
    print(f"  {match}  {label}")
    print(f"       v1: {s1}   v2: {s2}")
print()

if all_ok:
    print("  ✓ All expected edges present in both outputs")
else:
    print("  ⚠ Some expected edges differ between outputs")
print()
print(f"  v1 output : {v1_path}")
print(f"  v2 output : {v2_path}")
print()
print("  To inspect manually:")
print(f"    python3 -m json.tool {v2_path} | less")
PY
