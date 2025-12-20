"use client"

import { useAuth } from '../contexts/AuthContext'
import { useRouter } from 'next/navigation'
import { useEffect } from 'react'
import type { ComponentType } from 'react'

export default function withAuth<P extends object>(WrappedComponent: ComponentType<P>) {
  const WithAuthComponent = (props: P) => {
    const { user, loading } = useAuth()
    const router = useRouter()

    useEffect(() => {
      if (!loading && !user) {
        router.replace('/')
      }
    }, [user, loading, router])

    if (loading) {
      return <div>Loading...</div> // Or a spinner component
    }

    if (!user) {
      return null // Or a redirect component
    }

    return <WrappedComponent {...props} />
  }

  // Set a display name for easier debugging
  WithAuthComponent.displayName = `WithAuth(${WrappedComponent.displayName || WrappedComponent.name || 'Component'})`

  return WithAuthComponent
}
