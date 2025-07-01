import { createRouter, createWebHashHistory } from "vue-router";
import LoginView from "../views/LoginView.vue";
import ChatView from "../views/ChatView.vue";
import ProfileInfoView from "../views/ProfileInfoView.vue";
import AuthService from "../services/auth.js";

const routes = [
  {
    path: "/",
    redirect: "/chat",
  },
  {
    path: "/login",
    name: "Login",
    component: LoginView,
  },
  {
    path: "/chat",
    name: "Chat",
    component: ChatView,
  },
  {
    path: "/chat/:conversationId",
    name: "ChatConversation",
    component: ChatView,
    props: true,
  },
  {
    path: "/profile",
    name: "ProfileInfo",
    component: ProfileInfoView,
    props: (route) => ({
      type: route.query.type || "user",
      id: route.query.id || "me",
    }),
  },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

// Authentication guard
router.beforeEach((to, from, next) => {
  // Check if user is authenticated
  const isAuthenticated = AuthService.isAuthenticated();

  // If going to login page, allow
  if (to.name === "Login") {
    next();
    return;
  }

  // If not authenticated and not going to login, redirect to login
  if (!isAuthenticated) {
    next({ name: "Login" });
    return;
  }

  // If authenticated, allow navigation
  next();
});

export default router;
