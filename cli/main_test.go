/*
ビルドされたCLIコマンドのテストを行う。
いろんなコマンドを実行し、予想通りの結果が帰ってくるか。
これは各プログラムサーバ、APIGWサーバを複数起動した状態でテストする。

./cli -name toZip -i 6mb.txt -o out -p "10" -j
./cli -name toJson -i 6mb.txt -o out -p "10" -j
./cli -name err -i 6mb.txt -o out -p "10" -j
./cli -name sleep -i 6mb.txt -o out -p "10" -j

jsonが帰ってくるか
エラーならmsgsの中のメッセージは含まれているか
*/

package main_test
