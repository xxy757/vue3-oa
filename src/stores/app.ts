import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useAppStore = defineStore('app', () => {
  const collapsed = ref(false)
  const loading = ref(false)
  const breadcrumbs = ref<{ title: string; path?: string }[]>([])

  const isCollapsed = computed(() => collapsed.value)
  const isLoading = computed(() => loading.value)

  function toggleCollapsed(): void {
    collapsed.value = !collapsed.value
  }

  function setCollapsed(value: boolean): void {
    collapsed.value = value
  }

  function setLoading(value: boolean): void {
    loading.value = value
  }

  function setBreadcrumbs(breadcrumbList: { title: string; path?: string }[]): void {
    breadcrumbs.value = breadcrumbList
  }

  return {
    collapsed,
    loading,
    breadcrumbs,
    isCollapsed,
    isLoading,
    toggleCollapsed,
    setCollapsed,
    setLoading,
    setBreadcrumbs
  }
})
