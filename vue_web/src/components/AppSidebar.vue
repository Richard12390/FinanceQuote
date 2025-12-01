<script setup lang="ts">
import { RouterLink, useRouter } from "vue-router"
import { Home, Clock9, ClipboardClock, Search, Settings, Star } from "lucide-vue-next"
import { logout } from '@/api/auth';
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'

const items = [
  {
    title: "Home",
    url: "/index",
    icon: Home,
  },
  {
    title: "Real-time Quote",
    url: "/index",
    icon: Clock9,
  },
  {
    title: "Daily Quote",
    url: "/index",
    icon: ClipboardClock,
  },
  {
    title: "Search",
    url: "/index",
    icon: Search,
  },
  {
    title: "Settings",
    url: "/index",
    icon: Settings,
  },
  {
    title: "Favorites",
    url: "/favorites",
    icon: Star,
  },  
];

const router = useRouter();

const handleSignOut = async () => {
  try {
    await logout();
  } catch (error) {
    console.error("Logout failed:", error);
  } finally {
    localStorage.removeItem('loginUser');
    router.push('/login');
}
}
</script>


<template>
    <Sidebar variant="inset">
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupLabel>Menu</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
                <SidebarMenuItem v-for="item in items" :key="item.title">
                  <SidebarMenuButton asChild>
                      <RouterLink :to="item.url">
                        <component :is="item.icon" />
                        <span>{{item.title}}</span>
                      </RouterLink>
                  </SidebarMenuButton>                
                </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
      
    <SidebarFooter>
      <SidebarMenu>
        <SidebarMenuItem>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <SidebarMenuButton>
                <LucideUser2 /> User
                <LucideChevronUp class="ml-auto" />
              </SidebarMenuButton>
            </DropdownMenuTrigger>
            <DropdownMenuContent
              side="top"
              class="w-[--reka-popper-anchor-width] min-w-[11rem]" 
            >
              <DropdownMenuItem>
                <span>Account</span>
              </DropdownMenuItem>
              <DropdownMenuItem @click="handleSignOut">
                <span>Sign out</span>
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarFooter>
  </Sidebar>
</template>
