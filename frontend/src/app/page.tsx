'use client';

import { useEffect } from 'react';
import { getAccessToken } from '@/lib/auth';

export default function HomePage() {
  useEffect(() => {
    const token = getAccessToken();
    window.location.href = token ? '/dashboard' : '/login';
  }, []);

  return null;
}
