'use client';

import { useEffect } from 'react';
import i18n, { getSavedLocale } from '@/i18n/config';

export function I18nProvider({ children }: { children: React.ReactNode }) {
  useEffect(() => {
    const locale = getSavedLocale();
    if (locale && locale !== i18n.language) {
      i18n.changeLanguage(locale);
    }
  }, []);

  return <>{children}</>;
}
