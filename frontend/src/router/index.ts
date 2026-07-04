import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import MainLayout from '@/layouts/MainLayout.vue'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
  },
  {
    path: '/',
    component: MainLayout,
    redirect: '/consignment/consignors',
    children: [
      {
        path: 'consignment/consignors',
        name: 'ConsignorList',
        component: () => import('@/views/consignment/ConsignorList.vue'),
        meta: { title: '寄卖人管理' },
      },
      {
        path: 'consignment/items',
        name: 'ConsignmentList',
        component: () => import('@/views/consignment/ConsignmentList.vue'),
        meta: { title: '寄卖品管理' },
      },
      {
        path: 'consignment/settlements',
        name: 'SettlementList',
        component: () => import('@/views/consignment/SettlementList.vue'),
        meta: { title: '结算管理' },
      },
      {
        path: 'consignment/transfer',
        name: 'TransferProgress',
        component: () => import('@/views/consignment/TransferProgress.vue'),
        meta: { title: '过户进度' },
      },
      {
        path: 'store/products',
        name: 'ProductList',
        component: () => import('@/views/store/ProductList.vue'),
        meta: { title: '商品管理' },
      },
      {
        path: 'store/suppliers',
        name: 'SupplierList',
        component: () => import('@/views/store/SupplierList.vue'),
        meta: { title: '供应商管理' },
      },
      {
        path: 'store/purchases',
        name: 'PurchaseList',
        component: () => import('@/views/store/PurchaseList.vue'),
        meta: { title: '进货管理' },
      },
      {
        path: 'store/pos',
        name: 'PosCheckout',
        component: () => import('@/views/store/PosCheckout.vue'),
        meta: { title: 'POS收银' },
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem('token')
  if (to.path !== '/login' && !token) {
    next('/login')
  } else if (to.path === '/login' && token) {
    next('/')
  } else {
    next()
  }
})

export default router
