import { createRouter, createWebHashHistory } from "vue-router";
import HomeView from "../views/HomeView.vue";
import LoginView from "../views/LoginView.vue";
import ChatView from "../views/ChatView.vue";
import ProfileView from "../views/ProfileView.vue";

const routes = [
  {
    path: "/",
    name: "Home",
    component: HomeView,
  },
  {
    path: "/login",
    name: "Login",
    component: LoginView,
  },
  {
    path: "/chat/:id",
    name: "Chat",
    component: ChatView,
  },
  {
    path: "/profile",
    name: "Profile",
    component: ProfileView,
  },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

export default router;
