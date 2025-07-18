# logger

使用`NewLoggerWithWG`获取一个`Logger`。该`Logger`会定时刷新日志。

## cobra支持

使用`NewLogOptions`创建一个带有服务名称的日志配置。  
使用`Options.FlagSet`将log需要的flag设置到command中