package msglog

// 使用github.com/davyxu/golog的cellnet配色方案
const LogColorDefine = `
{
	"Rule":[
		{"Text":"panic:","Color":"Red"},
		{"Text":"[DB]","Color":"Green"},
		{"Text":"#http.listen","Color":"Blue"},
		{"Text":"#http.recv","Color":"Blue"},
		{"Text":"#http.send","Color":"Purple"},

		{"Text":"#tcp.listen","Color":"Blue"},
		{"Text":"#tcp.accepted","Color":"Blue"},
		{"Text":"#tcp.closed","Color":"Blue"},
		{"Text":"#tcp.recv","Color":"Blue"},
		{"Text":"#tcp.send","Color":"Purple"},
		{"Text":"#tcp.connected","Color":"Blue"},

		{"Text":"#ws.listen","Color":"Blue"},
		{"Text":"#ws.accepted","Color":"Blue"},
		{"Text":"#ws.closed","Color":"Blue"},
		{"Text":"#ws.recv","Color":"Blue"},
		{"Text":"#ws.send","Color":"Purple"},
		{"Text":"#ws.connected","Color":"Blue"},

		{"Text":"#udp.listen","Color":"Blue"},
		{"Text":"#udp.recv","Color":"Blue"},
		{"Text":"#udp.send","Color":"Purple"},

		{"Text":"#rpc.recv","Color":"Blue"},
		{"Text":"#rpc.send","Color":"Purple"},

		{"Text":"#relay.recv","Color":"Blue"},
		{"Text":"#relay.send","Color":"Purple"},

		{"Text":"#agent.recv","Color":"Blue"},
		{"Text":"#agent.send","Color":"Purple"}
	]
}
`
