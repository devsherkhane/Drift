<template>
  <div class="landing-page">
    <!-- Navigation -->
    <nav class="landing-nav">
      <div class="nav-brand">
        <div class="nav-logo-box">
          <ColumnsIcon :size="24" />
        </div>
        <span class="nav-logo-text">Trello Clone</span>
      </div>
      <div class="nav-actions">
        <router-link to="/login" class="nav-link">Sign In</router-link>
        <router-link to="/signup" class="nav-btn-primary">Get Started Free</router-link>
      </div>
    </nav>

    <!-- Hero Section -->
    <section class="hero-section">
      <div class="hero-blobs">
        <div class="hero-blob blob-a"></div>
        <div class="hero-blob blob-b"></div>
        <div class="hero-blob blob-c"></div>
      </div>
      <div class="hero-content animate-slide-up">
        <div class="hero-badge">
          <SparklesIcon :size="14" />
          <span>The Smarter Way to Manage Projects</span>
        </div>
        <h1 class="hero-title">
          Organize anything,<br />
          <span class="gradient-text">together.</span>
        </h1>
        <p class="hero-subtitle">
          Trello Clone brings all your tasks, teammates, and tools together.
          Keep everything in the same place — even if your team isn't.
        </p>
        <div class="hero-cta">
          <router-link to="/signup" class="btn-hero-primary">
            <RocketIcon :size="20" />
            Start for Free
          </router-link>
          <router-link to="/login" class="btn-hero-secondary">
            <LogInIcon :size="18" />
            Sign In
          </router-link>
        </div>
        <div class="hero-stats">
          <div class="stat-item">
            <span class="stat-number">10K+</span>
            <span class="stat-label">Active Users</span>
          </div>
          <div class="stat-divider"></div>
          <div class="stat-item">
            <span class="stat-number">50K+</span>
            <span class="stat-label">Boards Created</span>
          </div>
          <div class="stat-divider"></div>
          <div class="stat-item">
            <span class="stat-number">99.9%</span>
            <span class="stat-label">Uptime</span>
          </div>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section class="features-section">
      <div class="features-header animate-slide-up">
        <h2>Everything you need to <span class="gradient-text">stay productive</span></h2>
        <p>Simple, flexible, and powerful. The tools your team needs to move fast.</p>
      </div>
      <div class="features-grid">
        <div class="feature-card glass-panel animate-slide-up" v-for="(feature, i) in features" :key="i" :style="{ animationDelay: `${i * 0.1}s` }">
          <div class="feature-icon-box" :style="{ background: feature.bg }">
            <component :is="feature.icon" :size="24" :style="{ color: feature.color }" />
          </div>
          <h3>{{ feature.title }}</h3>
          <p>{{ feature.description }}</p>
        </div>
      </div>
    </section>

    <!-- CTA Section -->
    <section class="cta-section">
      <div class="cta-card glass-panel">
        <div class="cta-blobs">
          <div class="cta-blob cta-blob-1"></div>
          <div class="cta-blob cta-blob-2"></div>
        </div>
        <div class="cta-content">
          <h2>Ready to get started?</h2>
          <p>Join thousands of teams already using Trello Clone to ship products faster.</p>
          <router-link to="/signup" class="btn-hero-primary">
            <ZapIcon :size="20" />
            Create Your Free Account
          </router-link>
        </div>
      </div>
    </section>

    <!-- Footer -->
    <footer class="landing-footer">
      <div class="footer-content">
        <div class="footer-brand">
          <ColumnsIcon :size="20" class="footer-logo-icon" />
          <span>Trello Clone</span>
        </div>
        <p class="footer-copy">&copy; 2026 Trello Clone. Built with ❤️ as a learning project.</p>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import {
  Columns as ColumnsIcon,
  Sparkles as SparklesIcon,
  Rocket as RocketIcon,
  LogIn as LogInIcon,
  Zap as ZapIcon,
  LayoutDashboard as LayoutDashboardIcon,
  Users as UsersIcon,
  Bell as BellIcon,
  Search as SearchIcon,
  Shield as ShieldIcon,
  Palette as PaletteIcon
} from 'lucide-vue-next';

const auth = useAuthStore();
const router = useRouter();

onMounted(() => {
  if (auth.isAuthenticated) {
    router.replace('/');
  }
});

const features = [
  {
    icon: LayoutDashboardIcon,
    title: 'Intuitive Boards',
    description: 'Organize projects into boards with lists and cards. Drag-and-drop to rearrange in real-time.',
    bg: '#e0e7ff',
    color: '#6366f1'
  },
  {
    icon: UsersIcon,
    title: 'Team Collaboration',
    description: 'Invite teammates, assign tasks, and share boards. Work together seamlessly.',
    bg: '#d1fae5',
    color: '#10b981'
  },
  {
    icon: BellIcon,
    title: 'Real-time Updates',
    description: 'Get instant notifications when cards move, comments are added, or due dates approach.',
    bg: '#fef3c7',
    color: '#f59e0b'
  },
  {
    icon: SearchIcon,
    title: 'Powerful Search',
    description: 'Find any card across all your boards instantly with our global search feature.',
    bg: '#fce7f3',
    color: '#ec4899'
  },
  {
    icon: ShieldIcon,
    title: 'Secure by Default',
    description: 'JWT authentication, rate limiting, and security headers keep your data safe.',
    bg: '#e0e7ff',
    color: '#6366f1'
  },
  {
    icon: PaletteIcon,
    title: 'Dark & Light Themes',
    description: 'Switch between beautiful light and dark modes to match your preference.',
    bg: '#f3e8ff',
    color: '#8b5cf6'
  }
];
</script>

<style scoped>
.landing-page {
  min-height: 100vh;
  background: var(--bg-gradient);
  overflow-x: hidden;
}

/* Nav */
.landing-nav {
  padding: 16px 48px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  position: sticky;
  top: 0;
  z-index: 100;
  background: var(--surface-glass);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-bottom: 1px solid var(--border-subtle);
}

.nav-brand {
  display: flex;
  align-items: center;
  gap: 12px;
}

.nav-logo-box {
  width: 40px;
  height: 40px;
  background: var(--brand-primary);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.nav-logo-text {
  font-size: 20px;
  font-weight: 900;
  color: var(--text-main);
  letter-spacing: -1px;
}

.nav-actions {
  display: flex;
  align-items: center;
  gap: 16px;
}

.nav-link {
  color: var(--text-muted);
  text-decoration: none;
  font-weight: 700;
  font-size: 14px;
  padding: 10px 20px;
  border-radius: 12px;
  transition: all 0.2s;
}

.nav-link:hover {
  background: var(--surface-secondary);
  color: var(--text-main);
}

.nav-btn-primary {
  background: var(--brand-primary);
  color: white;
  text-decoration: none;
  font-weight: 700;
  font-size: 14px;
  padding: 10px 24px;
  border-radius: 12px;
  transition: all 0.25s;
  box-shadow: 0 4px 14px rgba(99, 102, 241, 0.3);
}

.nav-btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(99, 102, 241, 0.4);
}

/* Hero */
.hero-section {
  position: relative;
  padding: 100px 48px 80px;
  text-align: center;
  overflow: hidden;
}

.hero-blobs {
  position: absolute;
  inset: 0;
  z-index: 0;
  pointer-events: none;
}

.hero-blob {
  position: absolute;
  border-radius: 50%;
  filter: blur(100px);
  opacity: 0.12;
}

.blob-a {
  width: 500px;
  height: 500px;
  background: var(--brand-primary);
  top: -100px;
  left: -100px;
  animation: float 12s infinite alternate;
}

.blob-b {
  width: 400px;
  height: 400px;
  background: #10b981;
  bottom: -80px;
  right: -80px;
  animation: float 10s infinite alternate-reverse;
}

.blob-c {
  width: 300px;
  height: 300px;
  background: #f59e0b;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  animation: float 14s infinite alternate;
}

@keyframes float {
  from { transform: translate(0, 0) scale(1); }
  to { transform: translate(60px, 40px) scale(1.15); }
}

.hero-content {
  position: relative;
  z-index: 10;
  max-width: 750px;
  margin: 0 auto;
}

.hero-badge {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  background: var(--brand-primary-light);
  color: var(--brand-primary);
  padding: 8px 20px;
  border-radius: 50px;
  font-size: 13px;
  font-weight: 700;
  margin-bottom: 32px;
}

.hero-title {
  font-size: 64px;
  font-weight: 900;
  line-height: 1.1;
  letter-spacing: -3px;
  color: var(--text-main);
  margin-bottom: 24px;
}

.gradient-text {
  background: linear-gradient(135deg, var(--brand-primary), #8b5cf6, #ec4899);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.hero-subtitle {
  font-size: 18px;
  color: var(--text-muted);
  line-height: 1.7;
  max-width: 540px;
  margin: 0 auto 40px;
}

.hero-cta {
  display: flex;
  justify-content: center;
  gap: 16px;
  margin-bottom: 56px;
}

.btn-hero-primary {
  background: var(--brand-primary);
  color: white;
  text-decoration: none;
  padding: 16px 32px;
  border-radius: 16px;
  font-weight: 800;
  font-size: 16px;
  display: inline-flex;
  align-items: center;
  gap: 10px;
  transition: all 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  box-shadow: 0 6px 24px rgba(99, 102, 241, 0.3);
  border: none;
  cursor: pointer;
}

.btn-hero-primary:hover {
  transform: translateY(-3px);
  box-shadow: 0 12px 32px rgba(99, 102, 241, 0.4);
}

.btn-hero-secondary {
  background: var(--surface-primary);
  color: var(--text-main);
  text-decoration: none;
  padding: 16px 32px;
  border-radius: 16px;
  font-weight: 800;
  font-size: 16px;
  display: inline-flex;
  align-items: center;
  gap: 10px;
  transition: all 0.3s;
  border: 1px solid var(--border-subtle);
  box-shadow: var(--shadow-sm);
}

.btn-hero-secondary:hover {
  transform: translateY(-3px);
  box-shadow: var(--shadow-md);
  border-color: var(--brand-primary);
  color: var(--brand-primary);
}

.hero-stats {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 32px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.stat-number {
  font-size: 28px;
  font-weight: 900;
  color: var(--text-main);
  letter-spacing: -1px;
}

.stat-label {
  font-size: 13px;
  color: var(--text-muted);
  font-weight: 600;
}

.stat-divider {
  width: 1px;
  height: 40px;
  background: var(--border-subtle);
}

/* Features */
.features-section {
  padding: 80px 48px;
}

.features-header {
  text-align: center;
  margin-bottom: 56px;
}

.features-header h2 {
  font-size: 40px;
  font-weight: 900;
  color: var(--text-main);
  letter-spacing: -2px;
  margin-bottom: 16px;
}

.features-header p {
  font-size: 18px;
  color: var(--text-muted);
  max-width: 480px;
  margin: 0 auto;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;
  max-width: 1100px;
  margin: 0 auto;
}

.feature-card {
  background: var(--surface-primary);
  padding: 36px;
  border-radius: 20px;
  border: 1px solid var(--border-subtle);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.feature-card:hover {
  transform: translateY(-8px);
  box-shadow: var(--shadow-lg);
  border-color: var(--brand-primary-light);
}

.feature-icon-box {
  width: 52px;
  height: 52px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 20px;
}

.feature-card h3 {
  font-size: 18px;
  font-weight: 800;
  color: var(--text-main);
  margin-bottom: 10px;
  letter-spacing: -0.5px;
}

.feature-card p {
  font-size: 14px;
  color: var(--text-muted);
  line-height: 1.7;
}

/* CTA Section */
.cta-section {
  padding: 40px 48px 80px;
}

.cta-card {
  max-width: 900px;
  margin: 0 auto;
  padding: 72px 48px;
  text-align: center;
  border-radius: 28px;
  background: var(--surface-primary);
  position: relative;
  overflow: hidden;
}

.cta-blobs {
  position: absolute;
  inset: 0;
  pointer-events: none;
}

.cta-blob {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.08;
}

.cta-blob-1 {
  width: 300px;
  height: 300px;
  background: var(--brand-primary);
  top: -80px;
  right: -80px;
}

.cta-blob-2 {
  width: 250px;
  height: 250px;
  background: #10b981;
  bottom: -60px;
  left: -60px;
}

.cta-content {
  position: relative;
  z-index: 10;
}

.cta-content h2 {
  font-size: 36px;
  font-weight: 900;
  color: var(--text-main);
  letter-spacing: -2px;
  margin-bottom: 16px;
}

.cta-content p {
  font-size: 17px;
  color: var(--text-muted);
  margin-bottom: 32px;
  max-width: 440px;
  margin-left: auto;
  margin-right: auto;
  line-height: 1.6;
}

/* Footer */
.landing-footer {
  padding: 32px 48px;
  border-top: 1px solid var(--border-subtle);
}

.footer-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  max-width: 1100px;
  margin: 0 auto;
}

.footer-brand {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 16px;
  font-weight: 800;
  color: var(--text-main);
}

.footer-logo-icon {
  color: var(--brand-primary);
}

.footer-copy {
  font-size: 13px;
  color: var(--text-extramuted);
}

/* Animations */
.animate-slide-up {
  animation: slideUp 0.8s cubic-bezier(0.16, 1, 0.3, 1) both;
}

@keyframes slideUp {
  from { opacity: 0; transform: translateY(30px); }
  to { opacity: 1; transform: translateY(0); }
}

/* Responsive */
@media (max-width: 900px) {
  .features-grid { grid-template-columns: repeat(2, 1fr); }
  .hero-title { font-size: 44px; letter-spacing: -2px; }
  .landing-nav { padding: 16px 24px; }
  .hero-section { padding: 60px 24px; }
  .features-section { padding: 60px 24px; }
  .cta-section { padding: 40px 24px 60px; }
}

@media (max-width: 600px) {
  .features-grid { grid-template-columns: 1fr; }
  .hero-title { font-size: 36px; }
  .hero-cta { flex-direction: column; align-items: center; }
  .hero-stats { flex-direction: column; gap: 16px; }
  .stat-divider { width: 40px; height: 1px; }
  .footer-content { flex-direction: column; gap: 12px; text-align: center; }
}
</style>
