/**
 * 打印工具 - 使用新窗口打印，避免CSS冲突
 */

const SHOP_NAME = '家族门店'
const SHOP_PHONE = '请联系店员'

function openPrintWindow(title: string, bodyHtml: string) {
  const win = window.open('', '_blank', 'width=800,height=600')
  if (!win) {
    alert('请允许弹出窗口以进行打印')
    return
  }
  win.document.write(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<title>${title}</title>
<style>
  * { margin: 0; padding: 0; box-sizing: border-box; }
  body { font-family: "Microsoft YaHei", "PingFang SC", sans-serif; padding: 20px; color: #333; }
  .receipt { max-width: 380px; margin: 0 auto; }
  .receipt h2 { text-align: center; font-size: 18px; margin-bottom: 4px; }
  .receipt .subtitle { text-align: center; font-size: 12px; color: #666; margin-bottom: 12px; }
  .receipt hr { border: none; border-top: 1px dashed #999; margin: 8px 0; }
  .receipt .info-row { display: flex; justify-content: space-between; font-size: 13px; margin: 3px 0; }
  .receipt table { width: 100%; border-collapse: collapse; font-size: 13px; }
  .receipt th, .receipt td { padding: 4px 2px; text-align: left; }
  .receipt th { border-bottom: 1px solid #ddd; }
  .receipt .total { font-size: 16px; font-weight: bold; margin: 8px 0; text-align: right; }
  .receipt .footer { text-align: center; font-size: 11px; color: #999; margin-top: 12px; }
  .contract { max-width: 680px; margin: 0 auto; }
  .contract h2 { text-align: center; font-size: 20px; margin-bottom: 16px; }
  .contract h3 { font-size: 15px; margin: 16px 0 8px; }
  .contract p { font-size: 14px; line-height: 2; margin: 8px 0; }
  .contract .info-table { width: 100%; border-collapse: collapse; margin: 12px 0; font-size: 14px; }
  .contract .info-table td { padding: 6px 10px; border: 1px solid #ddd; }
  .contract .info-table .label { background: #f5f5f5; width: 120px; font-weight: 600; }
  .contract .sign-area { display: flex; justify-content: space-between; margin-top: 40px; font-size: 14px; }
  .contract .sign-item { flex: 1; }
  .contract .sign-line { border-top: 1px solid #333; margin-top: 40px; padding-top: 4px; text-align: center; }
  @media print { body { padding: 0; } }
</style>
</head>
<body>
${bodyHtml}
<script>window.onload=function(){window.print()}<\/script>
</body>
</html>`)
  win.document.close()
}

/** 打印POS小票 */
export function printReceipt(order: any, cart: any[], paymentName: string, totalAmount: number) {
  const now = new Date()
  const dateStr = now.toLocaleString('zh-CN')
  const orderNo = order?.id ? `#${order.id}` : `T${now.getTime()}`
  const itemsHtml = cart.map((item: any) => `
    <tr>
      <td>${item.name}</td>
      <td style="text-align:center">${item.qty}</td>
      <td style="text-align:right">¥${Number(item.price).toFixed(2)}</td>
      <td style="text-align:right">¥${(item.price * item.qty).toFixed(2)}</td>
    </tr>`).join('')

  const body = `<div class="receipt">
    <h2>${SHOP_NAME}</h2>
    <div class="subtitle">销售小票</div>
    <hr>
    <div class="info-row"><span>单号: ${orderNo}</span><span>${dateStr}</span></div>
    <div class="info-row"><span>收银员: ${order?.user_name || '管理员'}</span><span>顾客: 散客</span></div>
    <hr>
    <table>
      <tr><th>商品</th><th style="text-align:center">数量</th><th style="text-align:right">单价</th><th style="text-align:right">小计</th></tr>
      ${itemsHtml}
    </table>
    <hr>
    <div class="total">合计: ¥${totalAmount.toFixed(2)}</div>
    <div class="info-row"><span>支付方式: ${paymentName}</span><span>实付: ¥${totalAmount.toFixed(2)}</span></div>
    <hr>
    <div class="footer">感谢惠顾，欢迎再次光临！\n${SHOP_PHONE}</div>
  </div>`

  openPrintWindow('销售小票', body)
}

/** 打印寄卖合同 */
export function printConsignmentContract(item: any, consignor: any) {
  const dateStr = new Date().toLocaleDateString('zh-CN')
  const itemName = item.name || item.itemName || '-'
  const consignorName = consignor?.name || item.consignorName || item.consignor?.name || '-'
  const consignorPhone = consignor?.phone || item.consignor?.phone || '-'
  const expectedPrice = item.expected_price ?? item.expectedPrice ?? '-'
  const commission = item.commission_rate ?? item.commission ?? '-'
  const itemId = item.id || '-'

  const vehicleHtml = (item.is_vehicle || item.isVehicle) ? `
    <h3>车辆详情</h3>
    <table class="info-table">
      <tr><td class="label">品牌型号</td><td>${item.vehicle?.brand || '-'} ${item.vehicle?.model || ''}</td></tr>
      <tr><td class="label">车牌号</td><td>${item.vehicle?.plate_number || item.vehicle?.plateNumber || '-'}</td></tr>
      <tr><td class="label">车架号</td><td>${item.vehicle?.vin || '-'}</td></tr>
      <tr><td class="label">年份/颜色</td><td>${item.vehicle?.year || '-'}年 / ${item.vehicle?.color || '-'}</td></tr>
      <tr><td class="label">里程</td><td>${item.vehicle?.mileage ? item.vehicle.mileage + 'km' : '-'}</td></tr>
    </table>` : ''

  const body = `<div class="contract">
    <h2>寄卖协议</h2>
    <p>甲方（寄卖方）：${consignorName} &nbsp;&nbsp; 电话：${consignorPhone}</p>
    <p>乙方（代卖方）：${SHOP_NAME}</p>
    <p>签订日期：${dateStr}</p>

    <h3>寄卖物品信息</h3>
    <table class="info-table">
      <tr><td class="label">物品编号</td><td>#${itemId}</td></tr>
      <tr><td class="label">物品名称</td><td>${itemName}</td></tr>
      <tr><td class="label">物品描述</td><td>${item.description || item.memo || item.remark || '-'}</td></tr>
      <tr><td class="label">期望售价</td><td>${expectedPrice !== '-' ? '¥' + Number(expectedPrice).toLocaleString() : '-'}</td></tr>
      <tr><td class="label">佣金比例</td><td>${commission !== '-' ? commission + '%' : '面议'}</td></tr>
    </table>

    ${vehicleHtml}

    <h3>协议条款</h3>
    <p>1. 甲方委托乙方代为销售上述物品，销售价格以双方约定为准。</p>
    <p>2. 物品售出后，乙方按约定佣金比例收取服务费，余款支付给甲方。</p>
    <p>3. 寄卖期限内，甲方不得擅自取回物品。如需取回，应提前通知乙方。</p>
    <p>4. 物品在寄卖期间如有人为损坏，由责任方承担赔偿责任。</p>
    <p>5. 本协议一式两份，甲乙双方各执一份，自签订之日起生效。</p>

    <div class="sign-area">
      <div class="sign-item">
        <div>甲方签字：</div>
        <div class="sign-line">${consignorName}</div>
      </div>
      <div class="sign-item" style="margin-left:40px">
        <div>乙方签字：</div>
        <div class="sign-line">${SHOP_NAME}</div>
      </div>
    </div>
  </div>`

  openPrintWindow('寄卖协议', body)
}

/** 打印结算单 */
export function printSettlement(settlement: any, item: any, consignor: any) {
  const dateStr = new Date().toLocaleString('zh-CN')
  const itemName = item?.name || item?.itemName || '-'
  const consignorName = consignor?.name || item?.consignorName || item?.consignor?.name || '-'
  const typeMap: Record<string, string> = {
    SOLD_SETTLEMENT: '卖出结算',
    RETURN_SETTLEMENT: '到期取回',
    RENEWAL: '续费',
  }
  const typeName = typeMap[settlement.type] || settlement.type

  const detailRows = settlement.type === 'SOLD_SETTLEMENT' ? `
    <tr><td class="label">成交价格</td><td>¥${Number(settlement.sale_price || 0).toLocaleString()}</td></tr>
    <tr><td class="label">佣金</td><td>¥${Number(settlement.commission_amount || 0).toLocaleString()}</td></tr>
    <tr><td class="label">结算金额</td><td style="font-weight:bold;color:#e6a23c">¥${Number(settlement.settlement_amount || 0).toLocaleString()}</td></tr>` :
    settlement.type === 'RENEWAL' ? `
    <tr><td class="label">续费金额</td><td>¥${Number(settlement.renewal_fee || 0).toLocaleString()}</td></tr>
    <tr><td class="label">续费月数</td><td>${settlement.renewal_months || 1}个月</td></tr>` :
    `<tr><td class="label">结算类型</td><td>到期取回</td></tr>`

  const body = `<div class="contract">
    <h2>${typeName}单</h2>
    <table class="info-table">
      <tr><td class="label">结算编号</td><td>#${settlement.id || '-'}</td><td class="label">日期</td><td>${dateStr}</td></tr>
      <tr><td class="label">寄卖人</td><td>${consignorName}</td><td class="label">物品编号</td><td>#${item?.id || '-'}</td></tr>
      <tr><td class="label">物品名称</td><td colspan="3">${itemName}</td></tr>
    </table>

    <h3>结算明细</h3>
    <table class="info-table">
      ${detailRows}
      <tr><td class="label">备注</td><td>${settlement.remark || '-'}</td></tr>
    </table>

    <div class="sign-area">
      <div class="sign-item">
        <div>寄卖人签字：</div>
        <div class="sign-line">${consignorName}</div>
      </div>
      <div class="sign-item" style="margin-left:40px">
        <div>经办人签字：</div>
        <div class="sign-line">${SHOP_NAME}</div>
      </div>
    </div>
  </div>`

  openPrintWindow(`${typeName}单`, body)
}
