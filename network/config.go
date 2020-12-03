package network

import (
	"sync"
	"time"
)

type ServerInfo struct {
	Service
	Data   string
	Banner []string
}

type Service struct {
	AreaID int
	Id     int
	Type   int
	Name   string
}

type SessionConfig struct {
	Service
	Protocol          string        // 支持的协议 "tcp" "ws" "wss"
	Ip                string        // ip地址
	Port              int           // 端口
	IsClient          bool          // 客户端，链接发起方
	AutoReconnect     bool          // 是否自动重连
	ReconnectInterval time.Duration // 重连间隔
	IsInnerLink       bool          // 是否内部链接
	AuthKey           string        // Authentication Key
	SupportFragment   bool          // 是否支持分包发送
	ConnNum           int           // 客户端链接数量
	AllowMultiConn    bool          // 是否允许多链接
	Path              string        // ws
	CertFile          string        // wss Cert
	KeyFile           string        // wss Key
	MaxDone           int           // 接收队列缓存大小
	MaxSend           int           // 发送队列缓存大小
	MaxConn           int           // 最大链接数量
	MTU               int
	Linger            int
	NoDelay           bool
	KeepAlive         bool
	KeepAlivePeriod   time.Duration
	ReadBuffer        int
	WriteBuffer       int
	WriteTimeout      time.Duration
	ReadTimeout       time.Duration
	IdleTimeout       time.Duration // 多久未收到消息算空闲状态
	seq               int
	pkgDataPool       sync.Pool
	actionPool        sync.Pool

	//FilterChain []string
	//sfc         *SessionFilterChain
	//
	//HandlerChain []string
	//shc          *SessionHandlerChain
	//
	//ErrorPacketHandlerName string
	//eph                    ErrorPacketHandler
}

func (sc *SessionConfig) Init() {
	//if sc.EncoderName == "" {
	//	sc.encoder = packet.GetEncoder(packet.DefaultEncoderName)
	//} else {
	//	sc.encoder = packet.GetEncoder(sc.EncoderName)
	//}
	//if sc.DecoderName == "" {
	//	sc.decoder = packet.GetDecoder(packet.DefaultDecoderName)
	//} else {
	//	sc.decoder = packet.GetDecoder(sc.DecoderName)
	//}

	//for i := 0; i < len(sc.FilterChain); i++ {
	//	creator := GetSessionFilterCreator(sc.FilterChain[i])
	//	if creator != nil {
	//		if sc.sfc == nil {
	//			sc.sfc = NewSessionFilterChain()
	//		}
	//		if sc.sfc != nil {
	//			sc.sfc.AddLast(creator())
	//		}
	//	}
	//}
	//
	//for i := 0; i < len(sc.HandlerChain); i++ {
	//	creator := GetSessionHandlerCreator(sc.HandlerChain[i])
	//	if creator != nil {
	//		if sc.shc == nil {
	//			sc.shc = NewSessionHandlerChain()
	//		}
	//		if sc.shc != nil {
	//			sc.shc.AddLast(creator())
	//		}
	//	}
	//}
	//
	//if sc.ErrorPacketHandlerName != "" {
	//	creator := GetErrorPacketHandlerCreator(sc.ErrorPacketHandlerName)
	//	if creator != nil {
	//		sc.eph = creator()
	//	} else {
	//		logger.Logger.Warnf("[%v] ErrorPacketHandler not registe", sc.ErrorPacketHandlerName)
	//	}
	//}

	if sc.IdleTimeout <= 0 {
		sc.IdleTimeout = 5 * time.Second
	} else {
		sc.IdleTimeout *= time.Second
	}
	if sc.WriteTimeout <= 0 {
		sc.WriteTimeout = 30 * time.Second
	} else {
		sc.WriteTimeout *= time.Second
	}
	if sc.ReadTimeout <= 0 {
		sc.ReadTimeout = 30 * time.Second
	} else {
		sc.ReadTimeout *= time.Second
	}
	if sc.ReconnectInterval <= 0 {
		sc.ReconnectInterval = 5 * time.Second
	} else {
		sc.ReconnectInterval *= time.Second
	}
	sc.KeepAlivePeriod *= time.Second
	sc.pkgDataPool.New = func() interface{} {
		return new(PkgData)
	}
	sc.actionPool.New = func() interface{} {
		return new(action)
	}
}

func (sc *SessionConfig) GetSeq() int {
	sc.seq++
	return sc.seq
}

//func (sc *SessionConfig) GetFilter(name string) SessionFilter {
//	if sc.sfc != nil {
//		return sc.sfc.GetFilter(name)
//	}
//	return nil
//}
//
//func (sc *SessionConfig) GetHandler(name string) SessionHandler {
//	if sc.shc != nil {
//		return sc.shc.GetHandler(name)
//	}
//	return nil
//}
