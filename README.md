# kakscmd

Go library for serializing commands to Kakoune socket.

Import: `go.teddydd.me/kakscmd`.

```sh
make
./kak-raw-send -session 826387 -cmd 'eval -client client0 echo hello world' -debug
00000000  02 2f 00 00 00 26 00 00  00 65 76 61 6c 20 2d 63  |./...&...eval -c|
00000010  6c 69 65 6e 74 20 63 6c  69 65 6e 74 30 20 65 63  |lient client0 ec|
00000020  68 6f 20 68 65 6c 6c 6f  20 77 6f 72 6c 64 0a     |ho hello world.|

022F000000260000006576616C202D636C69656E7420636C69656E7430206563686F2068656C6C6F20776F726C640A
2022/05/03 14:17:28 written 47 bytes
```
