/**
 * Excel导出工具 - 使用原生方式生成CSV文件
 */

function downloadCsv(filename: string, headers: string[], rows: any[][]) {
  const csv = [headers.join(','), ...rows.map(row => row.map(cell => {
    const s = String(cell ?? '')
    return s.includes(',') || s.includes('"') || s.includes('\n') ? `"${s.replace(/"/g, '""')}"` : s
  }).join(','))].join('\n')
  
  const blob = new Blob(['\ufeff' + csv], { type: 'text/csv;charset=utf-8;' })
  const link = document.createElement('a')
  link.href = URL.createObjectURL(blob)
  link.download = filename
  link.click()
  URL.revokeObjectURL(link.href)
}

export function exportOrders(orders: any[]) {
  const headers = ['订单号', '客户名称', '支付方式', '金额', '状态', '创建时间']
  const rows = orders.map(o => [
    o.id || '-',
    o.customer_name || '-',
    o.payment_name || '-',
    o.total_paid ? `¥${Number(o.total_paid).toFixed(2)}` : '-',
    o.status || '-',
    o.created_at?.replace('T', ' ')?.substring(0, 19) || '-',
  ])
  downloadCsv(`订单记录_${new Date().toLocaleDateString('zh-CN').replace(/\//g, '-')}.csv`, headers, rows)
}

export function exportProducts(products: any[]) {
  const headers = ['商品ID', '名称', '分类', '价格', '库存', '单位', '供应商', '创建时间']
  const rows = products.map(p => [
    p.id || '-',
    p.name || '-',
    p.category_name || p.category?.name || '-',
    p.price ? `¥${Number(p.price).toFixed(2)}` : '-',
    p.stock || 0,
    p.base_unit || p.unit || '-',
    p.supplier_name || p.supplier?.name || '-',
    p.created_at?.replace('T', ' ')?.substring(0, 19) || '-',
  ])
  downloadCsv(`商品列表_${new Date().toLocaleDateString('zh-CN').replace(/\//g, '-')}.csv`, headers, rows)
}

const consignmentStatusMap: Record<string, string> = {
  ON_SALE: '在售',
  SOLD: '已售出',
  EXPIRED: '已过期',
  RETURNED: '已取回',
  CANCELLED: '已取消',
}

export function exportConsignments(items: any[]) {
  const headers = ['ID', '物品名称', '寄卖人', '状态', '期望价格', '成交价格', '佣金率', '是否车辆', '创建时间']
  const rows = items.map(item => [
    item.id || '-',
    item.name || item.itemName || '-',
    item.consignor_name || item.consignorName || item.consignor?.name || '-',
    item.status ? (consignmentStatusMap[String(item.status)] || item.status) : '-',
    (item.expected_price ?? item.expectedPrice) ? `¥${Number(item.expected_price ?? item.expectedPrice).toFixed(2)}` : '-',
    (item.final_price ?? item.sellingPrice) ? `¥${Number(item.final_price ?? item.sellingPrice).toFixed(2)}` : '-',
    (item.commission_rate ?? item.commission) ? `${item.commission_rate ?? item.commission}%` : '-',
    item.is_vehicle || item.isVehicle ? '是' : '否',
    (item.created_at || item.createdAt)?.replace('T', ' ')?.substring(0, 19) || '-',
  ])
  downloadCsv(`寄卖品记录_${new Date().toLocaleDateString('zh-CN').replace(/\//g, '-')}.csv`, headers, rows)
}

export function exportConsignors(consignors: any[]) {
  const headers = ['ID', '姓名', '电话', '身份证号', '地址', '创建时间']
  const rows = consignors.map(c => [
    c.id || '-',
    c.name || '-',
    c.phone || '-',
    c.id_card || '-',
    c.address || '-',
    c.created_at?.replace('T', ' ')?.substring(0, 19) || '-',
  ])
  downloadCsv(`寄卖人列表_${new Date().toLocaleDateString('zh-CN').replace(/\//g, '-')}.csv`, headers, rows)
}

export function exportSettlements(settlements: any[]) {
  const typeMap: Record<string, string> = {
    SOLD_SETTLEMENT: '卖出结算',
    RETURN_SETTLEMENT: '到期取回',
    RENEWAL: '续费',
  }
  const headers = ['ID', '类型', '成交价', '佣金', '结算金额', '续费金额', '续费月数', '备注', '时间']
  const rows = settlements.map(s => [
    s.id || '-',
    typeMap[s.type] || s.type,
    s.sale_price ? `¥${Number(s.sale_price).toFixed(2)}` : '-',
    s.commission_amount ? `¥${Number(s.commission_amount).toFixed(2)}` : '-',
    s.settlement_amount ? `¥${Number(s.settlement_amount).toFixed(2)}` : '-',
    s.renewal_fee ? `¥${Number(s.renewal_fee).toFixed(2)}` : '-',
    s.renewal_months || '-',
    s.remark || '-',
    s.created_at?.replace('T', ' ')?.substring(0, 19) || '-',
  ])
  downloadCsv(`结算记录_${new Date().toLocaleDateString('zh-CN').replace(/\//g, '-')}.csv`, headers, rows)
}

export function exportSuppliers(suppliers: any[]) {
  const headers = ['ID', '名称', '联系人', '电话', '地址', '备注', '创建时间']
  const rows = suppliers.map(s => [
    s.id || '-',
    s.name || '-',
    s.contact_name || '-',
    s.phone || '-',
    s.address || '-',
    s.remark || '-',
    s.created_at?.replace('T', ' ')?.substring(0, 19) || '-',
  ])
  downloadCsv(`供应商列表_${new Date().toLocaleDateString('zh-CN').replace(/\//g, '-')}.csv`, headers, rows)
}
