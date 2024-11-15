import type { MobileNavigations } from "$lib/types"
import Dashboard from "$lib/assets/icons/dashboard.png"
import Transaction from "$lib/assets/icons/transaction.png"

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