# -*- coding: utf-8 -*-
# Part 2: 章节四到十三 + 文档生成

# ===== 四、ECS宕机场景分类 =====
story.append(H1("四、ECS宕机场景分类"))
story.append(H2("4.1 场景分类总览"))
story.append(make_table(
    ["场景类型", "典型表现", "可能原因", "影响范围", "严重等级"],
    [
        ["单实例宕机", "单台ECS无法访问，控制台显示已停止", "宿主机故障、内核崩溃、OOM", "单个服务节点", "P3/P4"],
        ["批量实例宕机", "同一可用区多台ECS同时不可用", "宿主机集群故障、机房网络中断", "多个服务节点", "P2"],
        ["可用区级故障", "整个可用区ECS全部不可达", "机房电力故障、骨干网络中断", "整个可用区业务", "P1"],
        ["地域级故障", "整个地域所有可用区均受影响", "区域性自然灾害、大规模基础设施故障", "全部业务", "P1"],
    ],
    [25*mm, 35*mm, 40*mm, 30*mm, 20*mm]
))
story.append(SP())
story.append(H2("4.2 单实例宕机场景细分"))
story.append(make_table(
    ["子场景", "现象描述", "初步判断方法"],
    [
        ["实例停止", "控制台状态显示已停止，SSH无法连接", "检查控制台实例状态、查看系统事件"],
        ["实例无响应", "控制台显示运行中但SSH超时", "通过VNC登录检查、查看监控指标"],
        ["系统盘IO异常", "实例响应极慢，命令执行超时", "查看磁盘IO监控、检查iowait指标"],
        ["网络不可达", "Ping不通、端口不通但实例运行中", "检查安全组、VPC路由、EIP状态"],
        ["内存溢出OOM", "进程被Kill、服务异常退出", "查看dmesg日志、/var/log/messages"],
    ],
    [30*mm, 50*mm, 90*mm]
))
story.append(PageBreak())

# ===== 五、监控与告警机制 =====
story.append(H1("五、监控与告警机制"))
story.append(H2("5.1 云监控核心指标配置"))
story.append(make_table(
    ["监控指标", "告警阈值", "持续时间", "告警级别", "通知方式"],
    [
        ["CPU使用率", ">=90%", "连续5分钟", "警告", "钉钉/短信"],
        ["CPU使用率", ">=95%", "连续3分钟", "严重", "钉钉/短信/电话"],
        ["内存使用率", ">=85%", "连续5分钟", "警告", "钉钉/短信"],
        ["内存使用率", ">=95%", "连续3分钟", "严重", "钉钉/短信/电话"],
        ["磁盘使用率", ">=85%", "连续10分钟", "警告", "钉钉/短信"],
        ["磁盘IO等待", ">=80%", "连续5分钟", "严重", "钉钉/短信/电话"],
        ["系统状态检查", "失败", "连续1分钟", "紧急", "钉钉/短信/电话"],
        ["实例状态检查", "失败", "连续1分钟", "紧急", "钉钉/短信/电话"],
        ["网络入带宽", ">=带宽上限90%", "连续5分钟", "警告", "钉钉/短信"],
        ["TCP连接数", ">=上限80%", "连续5分钟", "警告", "钉钉/短信"],
    ],
    [30*mm, 28*mm, 25*mm, 20*mm, 32*mm]
))
story.append(SP())
story.append(H2("5.2 自定义监控项"))
story.append(BL("应用进程存活监控：关键进程（Nginx、Java、MySQL等）存活状态"))
story.append(BL("业务端口监控：核心业务端口（80/443/8080等）可达性检测"))
story.append(BL("日志关键字监控：ERROR、FATAL、OOM等关键字实时扫描"))
story.append(BL("自定义业务指标：QPS、响应时间、错误率等业务层面指标"))
story.append(SP())
story.append(H2("5.3 告警通知渠道"))
story.append(make_table(
    ["通知渠道", "适用场景", "配置要求", "响应时效"],
    [
        ["钉钉机器人", "所有级别告警", "配置Webhook地址，@相关人员", "实时"],
        ["短信通知", "警告及以上级别", "配置值班人员手机号码", "实时"],
        ["电话告警", "严重/紧急级别", "配置主备值班电话", "实时"],
        ["邮件通知", "所有级别（备份）", "配置运维组邮件列表", "5分钟内"],
        ["企业微信", "所有级别告警", "配置企业微信机器人Webhook", "实时"],
    ],
    [28*mm, 35*mm, 55*mm, 22*mm]
))
story.append(NOTE("建议至少配置两种以上通知渠道，避免单一渠道故障导致告警遗漏。"))
story.append(PageBreak())

# ===== 六、应急响应流程 =====
story.append(H1("六、应急响应流程"))
story.append(H2("6.1 故障分级标准"))
story.append(make_table(
    ["级别", "名称", "定义", "响应时间", "恢复目标RTO"],
    [
        ["P1", "紧急", "核心业务全面中断，影响全部用户或造成数据丢失风险", "5分钟", "30分钟"],
        ["P2", "严重", "核心业务部分中断，影响大量用户或关键功能不可用", "10分钟", "1小时"],
        ["P3", "一般", "非核心业务中断或核心业务性能严重下降", "30分钟", "4小时"],
        ["P4", "轻微", "非核心业务性能下降，影响少量用户", "1小时", "8小时"],
    ],
    [15*mm, 15*mm, 65*mm, 22*mm, 28*mm]
))
story.append(SP())
story.append(H2("6.2 应急响应流程总览"))
story.append(B("应急响应遵循以下六个阶段，形成闭环管理："))
story.append(SP())
flow_data = [
    [Paragraph("① 故障发现与确认", s_th), Paragraph("② 故障定级", s_th), Paragraph("③ 应急响应启动", s_th)],
    [Paragraph("监控告警触发或人工发现异常，值班人员5分钟内确认故障真实性", s_tc),
     Paragraph("根据影响范围和业务影响程度，确定P1-P4故障等级", s_tc),
     Paragraph("通知相关人员，建立应急沟通群，启动对应级别响应流程", s_tc)],
]
ft = Table(flow_data, colWidths=[56*mm, 56*mm, 56*mm])
ft.setStyle(TableStyle([
    ("BACKGROUND",(0,0),(-1,0),SECONDARY), ("BACKGROUND",(0,1),(-1,1),LIGHT_BG),
    ("ALIGN",(0,0),(-1,-1),"CENTER"), ("VALIGN",(0,0),(-1,-1),"MIDDLE"),
    ("GRID",(0,0),(-1,-1),0.5,BORDER_COLOR),
    ("TOPPADDING",(0,0),(-1,-1),3*mm), ("BOTTOMPADDING",(0,0),(-1,-1),3*mm),
]))
story.append(ft)
story.append(SP(2))
flow_data2 = [
    [Paragraph("④ 故障处置", s_th), Paragraph("⑤ 业务恢复", s_th), Paragraph("⑥ 故障复盘", s_th)],
    [Paragraph("技术处置组执行故障排查与修复，每15分钟更新进展", s_tc),
     Paragraph("验证服务恢复，逐步放量，持续观察30分钟确认稳定", s_tc),
     Paragraph("48小时内完成复盘报告，明确根因和改进措施", s_tc)],
]
ft2 = Table(flow_data2, colWidths=[56*mm, 56*mm, 56*mm])
ft2.setStyle(TableStyle([
    ("BACKGROUND",(0,0),(-1,0),SECONDARY), ("BACKGROUND",(0,1),(-1,1),LIGHT_BG),
    ("ALIGN",(0,0),(-1,-1),"CENTER"), ("VALIGN",(0,0),(-1,-1),"MIDDLE"),
    ("GRID",(0,0),(-1,-1),0.5,BORDER_COLOR),
    ("TOPPADDING",(0,0),(-1,-1),3*mm), ("BOTTOMPADDING",(0,0),(-1,-1),3*mm),
]))
story.append(ft2)
story.append(SP())
story.append(H2("6.3 故障发现与确认"))
story.append(H3("6.3.1 自动发现"))
story.append(BL("云监控告警：系统状态检查失败、实例状态检查失败自动触发"))
story.append(BL("自定义监控：应用进程监控、端口监控异常触发"))
story.append(BL("第三方监控：Zabbix、Prometheus、Grafana等外部监控系统告警"))
story.append(H3("6.3.2 人工发现"))
story.append(BL("用户反馈：客服或业务方报告服务异常"))
story.append(BL("巡检发现：日常巡检中发现实例异常"))
story.append(H3("6.3.3 确认步骤"))
story.append(BL("步骤1：登录阿里云控制台，检查实例运行状态"))
story.append(BL("步骤2：查看云监控指标，确认是否存在异常"))
story.append(BL("步骤3：尝试SSH/VNC连接实例，验证可达性"))
story.append(BL("步骤4：检查阿里云健康大盘，排除平台级故障"))
story.append(BL("步骤5：确认故障后，立即进入定级流程"))
story.append(SP())
story.append(H2("6.4 故障定级维度"))
story.append(make_table(
    ["评估维度", "P1紧急", "P2严重", "P3一般", "P4轻微"],
    [
        ["影响用户数", "全部用户", ">50%用户", "10%-50%用户", "<10%用户"],
        ["业务影响", "核心业务全面中断", "核心业务部分中断", "非核心业务中断", "性能轻微下降"],
        ["数据风险", "存在数据丢失风险", "数据写入受影响", "数据读取受影响", "无数据风险"],
        ["影响实例数", ">10台或整个集群", "5-10台", "2-4台", "1台"],
    ],
    [28*mm, 33*mm, 33*mm, 33*mm, 33*mm]
))
story.append(PageBreak())
