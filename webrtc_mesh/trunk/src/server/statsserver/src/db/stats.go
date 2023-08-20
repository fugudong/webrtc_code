package db

type CountryStat struct {
	Name 	string		// 国家名
	Count   int			// 计数
}

type BrowserStat struct {
	Name 	string		// 浏览器名
	Count   int			// 计数
}

type TalkStat struct {
	TalkTotal 			int // 通话总计
	AudioTalkTotal		int // 音频通话总计
	P2pRatio			float32 // p2p成功率
}