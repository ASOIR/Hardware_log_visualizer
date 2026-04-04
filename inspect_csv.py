import codecs

def detect_encoding(file_path):
    with open(file_path, 'rb') as f:
        raw = f.read(4)
    if raw.startswith(codecs.BOM_UTF8):
        return 'utf-8-sig'
    elif raw.startswith(codecs.BOM_UTF16_LE):
        return 'utf-16-le'
    elif raw.startswith(codecs.BOM_UTF16_BE):
        return 'utf-16-be'
    return 'utf-8'

encoding = detect_encoding('20260404-1.CSV')

with open('output.txt', 'w', encoding='utf-8') as out:
    out.write(f"Detected: {encoding}\n")
    with open('20260404-1.CSV', 'r', encoding=encoding, errors='replace') as f:
        for i in range(10):
            out.write(f.readline())
