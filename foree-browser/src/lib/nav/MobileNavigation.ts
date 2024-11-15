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
            {
                title: "Send Money",
                href: "/create_transaction",
                icon: Send
            }
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
            }
        ]
    }
]

export const sideMenuNavigations: Navigations = [
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
    {
        title: "Send Money",
        href: "/create_transaction",
        icon: Send
    }
]

export const headerDropdownNavigations: Navigations = [
    {
        title: "Dashboard",
        href: "/dashboard",
        icon: "not use"
    },
]