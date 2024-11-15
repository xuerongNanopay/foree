import type { MobileNavigations, Navigations } from "$lib/types"
import Dashboard from "$lib/assets/icons/dashboard.png"
import Transaction from "$lib/assets/icons/transaction.png"
import Send from "$lib/assets/icons/send.png"

export const mobileNavigations: MobileNavigations = [
    {
        subMenuTitle: "Menu",
        defaultActive: true,
        navigations: [
            {
                title: "Dashboard",
                href: "/dashboard",
                icon: Dashboard
            },
            {
                title: "Transaction",
                href: "/transaction",
                icon: Transaction
            },
        ]
    },
    {
        subMenuTitle: "Profile",
        navigations: [
            {
                title: "Personal",
                href: "/dashboard",
                icon: Dashboard
            },
            {
                title: "Notification Settings",
                href: "/transaction",
                icon: Transaction
            },
        ]
    }
]

export const sideMenuNavigations: Navigations = [
    {
        title: "Dashboard",
        href: "/dashboard",
        icon: "not use"
    },
    {
        title: "Transaction",
        href: "/transaction",
        icon: "not use"
    },
    {
        title: "Transfer",
        href: "/create_transaction",
        icon: "not use"
    }
]