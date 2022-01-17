package msgs

var (
	// SERVERERROR はサーバー内でエラーが起きた場合
	SERVERERROR string = "server error"

	// PROGRAMERROR は作成プログラムでエラーが起きた場合
	PROGRAMERROR string = "program error"

	// PROGRAMTIMEOUT は作成プログラムがタイムアウトした場合
	PROGRAMTIMEOUT string = "program timeout"

	// OK は上記全てのエラーが起きていない場合
	OK string = "ok"
)
