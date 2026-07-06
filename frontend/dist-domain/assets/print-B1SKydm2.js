const p="家族门店",u="请联系店员";function b(t,a){const e=window.open("","_blank","width=800,height=600");if(!e){alert("请允许弹出窗口以进行打印");return}e.document.write(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<title>${t}</title>
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
${a}
<script>window.onload=function(){window.print()}<\/script>
</body>
</html>`),e.document.close()}function S(t,a,e,i){const l=new Date,s=l.toLocaleString("zh-CN"),c=t!=null&&t.id?`#${t.id}`:`T${l.getTime()}`,n=a.map(d=>`
    <tr>
      <td>${d.name}</td>
      <td style="text-align:center">${d.qty}</td>
      <td style="text-align:right">¥${Number(d.price).toFixed(2)}</td>
      <td style="text-align:right">¥${(d.price*d.qty).toFixed(2)}</td>
    </tr>`).join(""),o=`<div class="receipt">
    <h2>${p}</h2>
    <div class="subtitle">销售小票</div>
    <hr>
    <div class="info-row"><span>单号: ${c}</span><span>${s}</span></div>
    <div class="info-row"><span>收银员: ${(t==null?void 0:t.user_name)||"管理员"}</span><span>顾客: 散客</span></div>
    <hr>
    <table>
      <tr><th>商品</th><th style="text-align:center">数量</th><th style="text-align:right">单价</th><th style="text-align:right">小计</th></tr>
      ${n}
    </table>
    <hr>
    <div class="total">合计: ¥${i.toFixed(2)}</div>
    <div class="info-row"><span>支付方式: ${e}</span><span>实付: ¥${i.toFixed(2)}</span></div>
    <hr>
    <div class="footer">感谢惠顾，欢迎再次光临！
${u}</div>
  </div>`;b("销售小票",o)}function E(t,a){var h,g,x,v,$,f,y,m,w,N;const e=new Date().toLocaleDateString("zh-CN"),i=t.name||t.itemName||"-",l=(a==null?void 0:a.name)||t.consignorName||((h=t.consignor)==null?void 0:h.name)||"-",s=(a==null?void 0:a.phone)||((g=t.consignor)==null?void 0:g.phone)||"-",c=t.expected_price??t.expectedPrice??"-",n=t.commission_rate??t.commission??"-",o=t.id||"-",d=t.is_vehicle||t.isVehicle?`
    <h3>车辆详情</h3>
    <table class="info-table">
      <tr><td class="label">品牌型号</td><td>${((x=t.vehicle)==null?void 0:x.brand)||"-"} ${((v=t.vehicle)==null?void 0:v.model)||""}</td></tr>
      <tr><td class="label">车牌号</td><td>${(($=t.vehicle)==null?void 0:$.plate_number)||((f=t.vehicle)==null?void 0:f.plateNumber)||"-"}</td></tr>
      <tr><td class="label">车架号</td><td>${((y=t.vehicle)==null?void 0:y.vin)||"-"}</td></tr>
      <tr><td class="label">年份/颜色</td><td>${((m=t.vehicle)==null?void 0:m.year)||"-"}年 / ${((w=t.vehicle)==null?void 0:w.color)||"-"}</td></tr>
      <tr><td class="label">里程</td><td>${(N=t.vehicle)!=null&&N.mileage?t.vehicle.mileage+"km":"-"}</td></tr>
    </table>`:"",r=`<div class="contract">
    <h2>寄卖协议</h2>
    <p>甲方（寄卖方）：${l} &nbsp;&nbsp; 电话：${s}</p>
    <p>乙方（代卖方）：${p}</p>
    <p>签订日期：${e}</p>

    <h3>寄卖物品信息</h3>
    <table class="info-table">
      <tr><td class="label">物品编号</td><td>#${o}</td></tr>
      <tr><td class="label">物品名称</td><td>${i}</td></tr>
      <tr><td class="label">物品描述</td><td>${t.description||t.memo||t.remark||"-"}</td></tr>
      <tr><td class="label">期望售价</td><td>${c!=="-"?"¥"+Number(c).toLocaleString():"-"}</td></tr>
      <tr><td class="label">佣金比例</td><td>${n!=="-"?n+"%":"面议"}</td></tr>
    </table>

    ${d}

    <h3>协议条款</h3>
    <p>1. 甲方委托乙方代为销售上述物品，销售价格以双方约定为准。</p>
    <p>2. 物品售出后，乙方按约定佣金比例收取服务费，余款支付给甲方。</p>
    <p>3. 寄卖期限内，甲方不得擅自取回物品。如需取回，应提前通知乙方。</p>
    <p>4. 物品在寄卖期间如有人为损坏，由责任方承担赔偿责任。</p>
    <p>5. 本协议一式两份，甲乙双方各执一份，自签订之日起生效。</p>

    <div class="sign-area">
      <div class="sign-item">
        <div>甲方签字：</div>
        <div class="sign-line">${l}</div>
      </div>
      <div class="sign-item" style="margin-left:40px">
        <div>乙方签字：</div>
        <div class="sign-line">${p}</div>
      </div>
    </div>
  </div>`;b("寄卖协议",r)}function z(t,a,e){var r;const i=new Date().toLocaleString("zh-CN"),l=(a==null?void 0:a.name)||(a==null?void 0:a.itemName)||"-",s=(e==null?void 0:e.name)||(a==null?void 0:a.consignorName)||((r=a==null?void 0:a.consignor)==null?void 0:r.name)||"-",n={SOLD_SETTLEMENT:"卖出结算",RETURN_SETTLEMENT:"到期取回",RENEWAL:"续费"}[t.type]||t.type,o=t.type==="SOLD_SETTLEMENT"?`
    <tr><td class="label">成交价格</td><td>¥${Number(t.sale_price||0).toLocaleString()}</td></tr>
    <tr><td class="label">佣金</td><td>¥${Number(t.commission_amount||0).toLocaleString()}</td></tr>
    <tr><td class="label">结算金额</td><td style="font-weight:bold;color:#e6a23c">¥${Number(t.settlement_amount||0).toLocaleString()}</td></tr>`:t.type==="RENEWAL"?`
    <tr><td class="label">续费金额</td><td>¥${Number(t.renewal_fee||0).toLocaleString()}</td></tr>
    <tr><td class="label">续费月数</td><td>${t.renewal_months||1}个月</td></tr>`:'<tr><td class="label">结算类型</td><td>到期取回</td></tr>',d=`<div class="contract">
    <h2>${n}单</h2>
    <table class="info-table">
      <tr><td class="label">结算编号</td><td>#${t.id||"-"}</td><td class="label">日期</td><td>${i}</td></tr>
      <tr><td class="label">寄卖人</td><td>${s}</td><td class="label">物品编号</td><td>#${(a==null?void 0:a.id)||"-"}</td></tr>
      <tr><td class="label">物品名称</td><td colspan="3">${l}</td></tr>
    </table>

    <h3>结算明细</h3>
    <table class="info-table">
      ${o}
      <tr><td class="label">备注</td><td>${t.remark||"-"}</td></tr>
    </table>

    <div class="sign-area">
      <div class="sign-item">
        <div>寄卖人签字：</div>
        <div class="sign-line">${s}</div>
      </div>
      <div class="sign-item" style="margin-left:40px">
        <div>经办人签字：</div>
        <div class="sign-line">${p}</div>
      </div>
    </div>
  </div>`;b(`${n}单`,d)}export{z as a,S as b,E as p};
