import re

def check_file(filepath):
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()

    # Find the Authorization table
    auth_table_match = re.search(r'## Authorization\n\n(.*?)\n\n---', content, re.DOTALL)
    if not auth_table_match:
        print(f"Could not find Authorization table in {filepath}")
        return

    auth_table = auth_table_match.group(1)
    
    # Extract endpoints from auth table (ignoring the header rows)
    table_endpoints = []
    for line in auth_table.strip().split('\n'):
        if line.startswith('| Endpoint') or line.startswith('| ---'):
            continue
        parts = line.split('|')
        if len(parts) > 1:
            ep = parts[1].strip()
            # sometimes there might be backticks or extra spaces
            ep = ep.replace('`', '').strip()
            table_endpoints.append(ep)

    # Find all detailed endpoints
    # e.g., **Endpoint:** `GET /api/v1/room-types`
    # or **Endpoint:** `GET /api/v1/rooms`
    detail_matches = re.findall(r'\*\*Endpoint:\*\*\s*`([A-Z]+)\s+(/api/v1/.*?)`', content)
    
    # Let's normalize the detail endpoints to match the table format.
    # The table format usually omits the `/api/v1` prefix.
    detail_endpoints = []
    for method, path in detail_matches:
        normalized_path = path.replace('/api/v1', '')
        if normalized_path == '':
            normalized_path = '/' # fallback though unlikely
        detail_endpoints.append(f"{method} {normalized_path}")

    print(f"=== Report for {filepath} ===")
    print(f"Endpoints in Authorization Table ({len(table_endpoints)}):")
    for ep in table_endpoints:
        status = "✅ Found in details" if ep in detail_endpoints else "❌ MISSING in details"
        print(f"  - {ep} -> {status}")
    
    print(f"\nEndpoints detailed below but NOT in Authorization Table:")
    extra_details = set(detail_endpoints) - set(table_endpoints)
    if extra_details:
        for ep in extra_details:
            print(f"  - {ep}")
    else:
        print("  None. All detailed endpoints are listed in the table.")
    print("\n")

check_file('06-room-types-api.md')
check_file('07-rooms-api.md')
