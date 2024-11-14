import type { MobileNavigation } from "$lib/types"
import Dashboard from "$lib/assets/icons/dashboard.png"
import Transaction from "$lib/assets/icons/transaction.png"

export const mobileNavigation: MobileNavigation = [
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
    }
]