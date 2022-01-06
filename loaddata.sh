#!/bin/bash

# Copy all data from files to database.

set -e

function dump_data {
    PGPASSWORD=postgres psql \
        --no-password -h localhost -U postgres -d stackoverflow_recommender -a \
        -c "\\copy ${1} FROM '${2}' DELIMITER '${3}' CSV HEADER;"
}

iconv -c -f UTF16 -t UTF8 data/U.csv -o data/U2.csv
dump_data "users" "data/U2.csv" ","
rm data/U2.csv

iconv -c -f UTF16 -t UTF8 data/Questions.csv -o data/Questions2.csv
dump_data "questions" "data/Questions2.csv" ","
rm data/Questions2.csv

cp data/tags.txt data/tags2.txt
echo "" >> data/tags2.txt
cat data/tags2.txt | awk -F " " '{print $1","$2}' >data/tags3.txt
sed 's/\r//' data/tags3.txt >data/tags4.txt
# iconv -c -f UTF16 -t UTF8 data/tags3.txt -o data/tags4.txt
dump_data "question_tags(question_id, tag)" "data/tags4.txt" ","
rm data/tags2.txt
rm data/tags3.txt
rm data/tags4.txt

iconv -c -f UTF16 -t UTF8 data/Answer.csv -o data/Answer2.csv
dump_data "answers" "data/Answer2.csv" ","
rm data/Answer2.csv

iconv -c -f UTF16 -t UTF8 data/Q_A.csv -o data/Q_A2.csv
dump_data "question_answers(question_id, answer_id, accepted)" "data/Q_A2.csv" ","
rm data/Q_A2.csv
