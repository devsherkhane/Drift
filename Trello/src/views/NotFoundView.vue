<template>
  <div class="not-found-page">
    <div class="nf-blobs">
      <div class="nf-blob nf-blob-1"></div>
      <div class="nf-blob nf-blob-2"></div>
    </div>

    <div class="nf-content animate-slide-up">
      <div class="nf-code">
        <span class="digit d-4-1">4</span>
        <div class="digit-zero">
          <div class="zero-orbit"></div>
          <span class="d-0">0</span>
        </div>
        <span class="digit d-4-2">4</span>
      </div>

      <h1 class="nf-title">Page not found</h1>
      <p class="nf-subtitle">
        Oops! The page you're looking for doesn't exist or has been moved.
      </p>

      <div class="nf-actions">
        <router-link to="/" class="btn-nf-primary">
          <HomeIcon :size="20" />
          Back to Dashboard
        </router-link>
        <router-link to="/search" class="btn-nf-secondary">
          <SearchIcon :size="18" />
          Search
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import {
  Home as HomeIcon,
  Search as SearchIcon
} from 'lucide-vue-next';
</script>

<style scoped>
.not-found-page {
  min-height: 100vh;
  background: var(--bg-gradient);
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  padding: 40px;
}

.nf-blobs {
  position: absolute;
  inset: 0;
  pointer-events: none;
}

.nf-blob {
  position: absolute;
  border-radius: 50%;
  filter: blur(100px);
  opacity: 0.1;
}

.nf-blob-1 {
  width: 500px;
  height: 500px;
  background: var(--brand-primary);
  top: -150px;
  right: -100px;
  animation: blobFloat 12s infinite alternate;
}

.nf-blob-2 {
  width: 400px;
  height: 400px;
  background: #ec4899;
  bottom: -100px;
  left: -100px;
  animation: blobFloat 10s infinite alternate-reverse;
}

@keyframes blobFloat {
  from { transform: translate(0, 0) scale(1); }
  to { transform: translate(40px, 30px) scale(1.1); }
}

.nf-content {
  text-align: center;
  position: relative;
  z-index: 10;
}

.nf-code {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin-bottom: 32px;
}

.digit {
  font-size: 120px;
  font-weight: 900;
  letter-spacing: -6px;
  background: linear-gradient(135deg, var(--brand-primary), #8b5cf6);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  line-height: 1;
}

.digit-zero {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}

.d-0 {
  font-size: 120px;
  font-weight: 900;
  letter-spacing: -6px;
  background: linear-gradient(135deg, #ec4899, #f59e0b);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  line-height: 1;
}

.zero-orbit {
  position: absolute;
  width: 90px;
  height: 90px;
  border: 3px dashed var(--brand-primary-light);
  border-radius: 50%;
  animation: orbitSpin 8s linear infinite;
}

@keyframes orbitSpin { to { transform: rotate(360deg); } }

.nf-title {
  font-size: 32px;
  font-weight: 900;
  color: var(--text-main);
  letter-spacing: -1.5px;
  margin-bottom: 12px;
}

.nf-subtitle {
  font-size: 17px;
  color: var(--text-muted);
  max-width: 400px;
  margin: 0 auto 40px;
  line-height: 1.6;
}

.nf-actions {
  display: flex;
  justify-content: center;
  gap: 16px;
}

.btn-nf-primary {
  background: var(--brand-primary);
  color: white;
  text-decoration: none;
  padding: 14px 28px;
  border-radius: 14px;
  font-weight: 800;
  font-size: 15px;
  display: inline-flex;
  align-items: center;
  gap: 10px;
  transition: all 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  box-shadow: 0 6px 20px rgba(99, 102, 241, 0.3);
}

.btn-nf-primary:hover {
  transform: translateY(-3px);
  box-shadow: 0 10px 28px rgba(99, 102, 241, 0.4);
}

.btn-nf-secondary {
  background: var(--surface-primary);
  color: var(--text-main);
  text-decoration: none;
  padding: 14px 28px;
  border-radius: 14px;
  font-weight: 800;
  font-size: 15px;
  display: inline-flex;
  align-items: center;
  gap: 10px;
  transition: all 0.3s;
  border: 1px solid var(--border-subtle);
  box-shadow: var(--shadow-sm);
}

.btn-nf-secondary:hover {
  transform: translateY(-3px);
  box-shadow: var(--shadow-md);
  border-color: var(--brand-primary);
  color: var(--brand-primary);
}

.animate-slide-up {
  animation: slideUp 0.8s cubic-bezier(0.16, 1, 0.3, 1) both;
}

@keyframes slideUp {
  from { opacity: 0; transform: translateY(30px); }
  to { opacity: 1; transform: translateY(0); }
}

@media (max-width: 600px) {
  .digit, .d-0 { font-size: 80px; }
  .zero-orbit { width: 60px; height: 60px; }
  .nf-actions { flex-direction: column; align-items: center; }
}
</style>
