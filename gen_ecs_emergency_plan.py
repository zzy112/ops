# -*- coding: utf-8 -*-
"""阿里云服务器ECS宕机应急预案 PDF生成脚本"""

from reportlab.lib.pagesizes import A4
from reportlab.lib.units import mm, cm
from reportlab.lib.styles import getSampleStyleSheet, ParagraphStyle
from reportlab.lib.colors import HexColor, black, white, grey
from reportlab.lib.enums import TA_CENTER, TA_LEFT, TA_JUSTIFY
from reportlab.platypus import (
    SimpleDocTemplate, Paragraph, Spacer, Table, TableStyle,
    PageBreak, KeepTogether
)
from reportlab.pdfbase import pdfmetrics
from reportlab.pdfbase.ttfonts import TTFont
import os

# ========== 字体注册 ==========
FONT_DIR = r"C:\Windows\Fonts"
pdfmetrics.registerFont(TTFont("SimHei", os.path.join(FONT_DIR, "simhei.ttf")))
pdfmetrics.registerFont(TTFont("SimSun", os.path.join(FONT_DIR, "simsun.ttc")))
pdfmetrics.registerFont(TTFont("MSYH", os.path.join(FONT_DIR, "msyh.ttc")))
pdfmetrics.registerFont(TTFont("MSYHBD", os.path.join(FONT_DIR, "msyhbd.ttc")))

# ========== 颜色 ==========
PRIMARY = HexColor("#1a5276")
SECONDARY = HexColor("#2e86c1")
LIGHT_BG = HexColor("#eaf2f8")
HEADER_BG = HexColor("#1a5276")
ROW_ALT = HexColor("#f5f8fa")
BORDER_COLOR = HexColor("#aab7c4")
ACCENT = HexColor("#e74c3c")

# ========== 样式 ==========
s_title = ParagraphStyle("CTitle", fontName="MSYHBD", fontSize=26, leading=36,
    alignment=TA_CENTER, textColor=PRIMARY, spaceAfter=6*mm)
s_subtitle = ParagraphStyle("CSubtitle", fontName="MSYH", fontSize=12, leading=18,
    alignment=TA_CENTER, textColor=HexColor("#555555"), spaceAfter=3*mm)
s_h1 = ParagraphStyle("CH1", fontName="MSYHBD", fontSize=16, leading=24,
    textColor=white, spaceBefore=8*mm, spaceAfter=4*mm, backColor=PRIMARY,
    leftIndent=4*mm, rightIndent=4*mm, borderPadding=(3*mm, 4*mm, 3*mm, 4*mm))
s_h2 = ParagraphStyle("CH2", fontName="MSYHBD", fontSize=13, leading=20,
    textColor=PRIMARY, spaceBefore=5*mm, spaceAfter=3*mm, leftIndent=2*mm)
s_h3 = ParagraphStyle("CH3", fontName="MSYHBD", fontSize=11, leading=17,
    textColor=SECONDARY, spaceBefore=3*mm, spaceAfter=2*mm, leftIndent=4*mm)
s_body = ParagraphStyle("CBody", fontName="MSYH", fontSize=10, leading=16,
    textColor=HexColor("#333333"), spaceAfter=2*mm, leftIndent=4*mm, alignment=TA_JUSTIFY)
s_body_bold = ParagraphStyle("CBodyBold", parent=s_body, fontName="MSYHBD")
s_bullet = ParagraphStyle("CBullet", parent=s_body, leftIndent=10*mm,
    bulletIndent=6*mm, spaceBefore=1*mm, spaceAfter=1*mm)
s_code = ParagraphStyle("CCode", fontName="SimSun", fontSize=9, leading=14,
    textColor=HexColor("#c0392b"), backColor=HexColor("#f9f2f4"),
    leftIndent=8*mm, rightIndent=4*mm, spaceBefore=2*mm, spaceAfter=2*mm,
    borderPadding=(2*mm, 3*mm, 2*mm, 3*mm))
s_note = ParagraphStyle("CNote", fontName="MSYH", fontSize=9, leading=14,
    textColor=HexColor("#7f8c8d"), leftIndent=8*mm, spaceBefore=1*mm,
    spaceAfter=2*mm, backColor=HexColor("#fef9e7"),
    borderPadding=(2*mm, 3*mm, 2*mm, 3*mm))
s_th = ParagraphStyle("CTH", fontName="MSYHBD", fontSize=9, leading=14,
    textColor=white, alignment=TA_CENTER)
s_tc = ParagraphStyle("CTC", fontName="MSYH", fontSize=9, leading=14,
    textColor=HexColor("#333333"), alignment=TA_CENTER)
s_tcl = ParagraphStyle("CTCL", fontName="MSYH", fontSize=9, leading=14,
    textColor=HexColor("#333333"))

# ========== 工具函数 ==========
def make_table(headers, rows, col_widths=None):
    hdr = [Paragraph(h, s_th) for h in headers]
    data = [hdr]
    for row in rows:
        data.append([Paragraph(str(c), s_tc if len(str(c)) < 20 else s_tcl) for c in row])
    w = col_widths or [170*mm // len(headers)] * len(headers)
    t = Table(data, colWidths=w, repeatRows=1)
    cmds = [
        ("BACKGROUND", (0, 0), (-1, 0), HEADER_BG),
        ("TEXTCOLOR", (0, 0), (-1, 0), white),
        ("ALIGN", (0, 0), (-1, -1), "CENTER"),
        ("VALIGN", (0, 0), (-1, -1), "MIDDLE"),
        ("GRID", (0, 0), (-1, -1), 0.5, BORDER_COLOR),
        ("TOPPADDING", (0, 0), (-1, -1), 3*mm),
        ("BOTTOMPADDING", (0, 0), (-1, -1), 3*mm),
        ("LEFTPADDING", (0, 0), (-1, -1), 2*mm),
        ("RIGHTPADDING", (0, 0), (-1, -1), 2*mm),
    ]
    for i in range(1, len(data)):
        if i % 2 == 0:
            cmds.append(("BACKGROUND", (0, i), (-1, i), ROW_ALT))
    t.setStyle(TableStyle(cmds))
    return t

def B(t): return Paragraph(t, s_body)
def BB(t): return Paragraph(t, s_body_bold)
def BL(t): return Paragraph(f"• {t}", s_bullet)
def H1(t): return Paragraph(t, s_h1)
def H2(t): return Paragraph(t, s_h2)
def H3(t): return Paragraph(t, s_h3)
def CODE(t): return Paragraph(t.replace("\n","<br/>"), s_code)
def NOTE(t): return Paragraph(f"⚠ {t}", s_note)
def SP(h=3): return Spacer(1, h*mm)

def header_footer(canvas, doc):
    canvas.saveState()
    canvas.setFont("MSYH", 8)
    canvas.setFillColor(HexColor("#999999"))
    canvas.drawString(20*mm, 290*mm, "阿里云ECS宕机应急预案 V1.0")
    canvas.drawRightString(190*mm, 290*mm, "内部机密")
    canvas.line(20*mm, 289*mm, 190*mm, 289*mm)
    canvas.drawCentredString(105*mm, 8*mm, f"- {doc.page} -")
    canvas.restoreState()

# ========== 构建文档 ==========
story = []

# ===== 封面 =====
story.append(Spacer(1, 60*mm))
story.append(Paragraph("阿里云服务器ECS宕机应急预案", s_title))
story.append(SP(6))
story.append(Paragraph("Emergency Response Plan for Alibaba Cloud ECS Downtime", s_subtitle))
story.append(SP(10))
cover_info = [
    ["文档编号", "OPS-ECS-ERP-2026-001"], ["版本号", "V1.0"], ["密级", "内部机密"],
    ["编制日期", "2026年04月02日"], ["编制部门", "运维保障部"],
    ["审核人", "____________"], ["批准人", "____________"],
]
cover_data = [[Paragraph(r[0], s_th), Paragraph(r[1], s_tc)] for r in cover_info]
ct = Table(cover_data, colWidths=[50*mm, 80*mm])
ct.setStyle(TableStyle([
    ("BACKGROUND", (0,0),(0,-1), HEADER_BG), ("TEXTCOLOR", (0,0),(0,-1), white),
    ("GRID", (0,0),(-1,-1), 0.5, BORDER_COLOR), ("VALIGN", (0,0),(-1,-1), "MIDDLE"),
    ("TOPPADDING", (0,0),(-1,-1), 3*mm), ("BOTTOMPADDING", (0,0),(-1,-1), 3*mm),
]))
story.append(ct)
story.append(SP(20))
story.append(Paragraph("本文档为内部机密文件，未经授权不得外传", ParagraphStyle(
    "warn", fontName="MSYHBD", fontSize=10, alignment=TA_CENTER, textColor=ACCENT)))
story.append(PageBreak())

# ===== 修订记录 =====
story.append(H1("修订记录"))
story.append(make_table(["版本","日期","修订人","修订内容","审批人"],
    [["V1.0","2026-04-02","运维保障部","初始版本发布","待审批"],["","","","",""]],
    [20*mm,28*mm,30*mm,60*mm,30*mm]))
story.append(PageBreak())

# ===== 目录 =====
story.append(H1("目  录"))
toc_items = ["一、目的与适用范围","二、术语定义","三、应急组织架构","四、ECS宕机场景分类",
    "五、监控与告警机制","六、应急响应流程","七、详细处置步骤","八、关键操作命令与API调用",
    "九、数据备份与恢复策略","十、业务切换方案","十一、通报与升级机制","十二、应急演练计划","十三、附录"]
for item in toc_items:
    story.append(Paragraph(item, ParagraphStyle("toc", fontName="MSYH",
        fontSize=11, leading=22, leftIndent=10*mm, textColor=PRIMARY)))
story.append(PageBreak())
