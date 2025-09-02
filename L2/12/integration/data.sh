#!/bin/bash

# Скрипт для создания и наполнения testdata файлов
# Usage: ./generate_testdata.sh

set -e  # Выход при первой ошибке

# Создаем директории testdata
mkdir -p testdata/{basic,case_insensitive,count_only,with_context,inverted,fixed_string,regex}

# Basic pattern matching
cat > testdata/basic/input.txt << 'EOF'
line 1
line with error
line 3
another error line
line 5
EOF

cat > testdata/basic/expected.txt << 'EOF'
line with error
another error line
EOF

# Case insensitive
cat > testdata/case_insensitive/input.txt << 'EOF'
line 1
line with error
LINE WITH ERROR
line 4
EOF

cat > testdata/case_insensitive/expected.txt << 'EOF'
line with error
LINE WITH ERROR
EOF

# Count only
cat > testdata/count_only/input.txt << 'EOF'
test line
another test
no match
TEST
EOF

cat > testdata/count_only/expected.txt << 'EOF'
3
EOF

# With context
cat > testdata/with_context/input.txt << 'EOF'
line 1
line 2
target line
line 4
line 5
another target
line 7
EOF

cat > testdata/with_context/expected.txt << 'EOF'
line 2
target line
line 4
another target
line 7
EOF

# Inverted match
cat > testdata/inverted/input.txt << 'EOF'
keep this
skip this
keep that
skip that
EOF

cat > testdata/inverted/expected.txt << 'EOF'
keep this
keep that
EOF

# Fixed string matching
cat > testdata/fixed_string/input.txt << 'EOF'
test line
[test] pattern
testing
no match
EOF

cat > testdata/fixed_string/expected.txt << 'EOF'
test line
[test] pattern
EOF

# Regex matching
cat > testdata/regex/input.txt << 'EOF'
test line
tast pattern
tost value
no match
EOF

cat > testdata/regex/expected.txt << 'EOF'
test line
tast pattern
tost value
EOF

echo "Testdata files created successfully!"
echo "Directory structure:"
find testdata -type f -name "*.txt" | sort