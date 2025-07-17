import { createRouter, createWebHistory } from 'vue-router'
import MainMenu from '../views/MainMenu.vue'
import CreateWallet from '../views/CreateWallet.vue'
import CheckWallet from '../views/CheckWallet.vue'
import WalletDetails from '../views/WalletDetails.vue'
import SendBitcoin from '../views/SendBitcoin.vue'

const routes = [
  {
    path: '/',
    name: 'MainMenu',
    component: MainMenu
  },
  {
    path: '/create-wallet',
    name: 'CreateWallet',
    component: CreateWallet
  },
  {
    path: '/check-wallet',
    name: 'CheckWallet',
    component: CheckWallet
  },
  {
    path: '/wallet-details',
    name: 'WalletDetails',
    component: WalletDetails
  },
  {
    path: '/send-bitcoin',
    name: 'SendBitcoin',
    component: SendBitcoin
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router