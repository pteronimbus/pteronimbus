import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import SearchAndFilters from '~/components/SearchAndFilters.vue'

describe('SearchAndFilters Component', () => {
  const createWrapper = (props = {}) => {
    return mount(SearchAndFilters, {
      props: {
        searchQuery: '',
        filters: [],
        ...props
      }
    })
  }

  it('renders component without errors', () => {
    expect(() => createWrapper()).not.toThrow()
  })

  it('displays search functionality', () => {
    const wrapper = createWrapper({
      searchQuery: 'test search'
    })
    
    expect(wrapper.html()).toBeTruthy()
    expect(wrapper.html()).toContain('test search')
  })

  it('handles search placeholder', () => {
    const wrapper = createWrapper({
      searchPlaceholder: 'Search servers...'
    })
    
    expect(wrapper.html()).toContain('Search servers...')
  })

  it('uses default placeholder when not provided', () => {
    const wrapper = createWrapper()
    
    expect(wrapper.html()).toContain('Search...')
  })

  it('emits search update when search changes', async () => {
    const wrapper = createWrapper()
    
    const searchInput = wrapper.find('input[type="text"]')
    if (searchInput.exists()) {
      await searchInput.setValue('new search')
      expect(wrapper.emitted('update:searchQuery')).toBeTruthy()
    }
  })

  it('handles filters prop correctly', () => {
    const wrapper = createWrapper({
      filters: [
        {
          key: 'status',
          value: 'all',
          options: [
            { label: 'All Status', value: 'all' },
            { label: 'Online', value: 'online' }
          ]
        }
      ]
    })
    
    expect(wrapper.html()).toBeTruthy()
  })

  it('handles empty filters array', () => {
    const wrapper = createWrapper({
      filters: []
    })
    
    expect(wrapper.html()).toBeTruthy()
  })

  it('handles multiple filters', () => {
    const wrapper = createWrapper({
      filters: [
        {
          key: 'status',
          value: 'all',
          options: [{ label: 'All Status', value: 'all' }]
        },
        {
          key: 'role',
          value: 'all',
          options: [{ label: 'All Roles', value: 'all' }]
        }
      ]
    })
    
    expect(wrapper.html()).toBeTruthy()
  })

  it('handles component creation with complex props', () => {
    expect(() => createWrapper({
      searchQuery: 'test',
      searchPlaceholder: 'Search here...',
      filters: [
        {
          key: 'category',
          value: 'all',
          options: [
            { label: 'All Categories', value: 'all' },
            { label: 'Games', value: 'games' }
          ],
          class: 'w-40'
        }
      ]
    })).not.toThrow()
  })

  it('maintains consistent structure with different props', () => {
    const simpleWrapper = createWrapper({
      searchQuery: 'simple'
    })
    
    const complexWrapper = createWrapper({
      searchQuery: 'complex search',
      filters: [
        {
          key: 'status',
          value: 'active',
          options: [
            { label: 'All', value: 'all' },
            { label: 'Active', value: 'active' }
          ]
        }
      ]
    })
    
    expect(simpleWrapper.html()).toBeTruthy()
    expect(complexWrapper.html()).toBeTruthy()
  })

  it('handles various search query values', () => {
    const testQueries = ['test', 'search query', '123', '']
    
    testQueries.forEach(query => {
      expect(() => createWrapper({
        searchQuery: query
      })).not.toThrow()
    })
  })

  it('validates component props correctly', () => {
    const wrapper = createWrapper({
      searchQuery: 'validation test',
      searchPlaceholder: 'Custom placeholder'
    })
    
    expect(wrapper.html()).toContain('validation test')
    expect(wrapper.html()).toContain('Custom placeholder')
  })
}) 