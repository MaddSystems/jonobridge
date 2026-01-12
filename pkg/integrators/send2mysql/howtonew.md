1. 🔄 Reemplazar todas las apariciones de “suntech” (incluso en palabras compuestas):
grep -irl "suntech" . | xargs sed -i 's/suntech/huabao/g'

2. 📁 Renombrar carpetas con “suntech” en el nombre:
find . -depth -type d -name '*suntech*' | while read dir; do mv "$dir" "$(echo "$dir" | sed 's/suntech/huabao/g')"; done

3. 📄 (Opcional) Renombrar archivos que contengan “suntech” en su nombre:
find . -type f -name '*suntech*' | while read file; do mv "$file" "$(echo "$file" | sed 's/suntech/huabao/g')"; done
