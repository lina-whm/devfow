import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import en from './en.json';
import ru from './ru.json';

const resources = { en: { translation: en }, ru: { translation: ru } };

i18n.use(initReactI18next).init({
  resources,
  lng: 'en',
  fallbackLng: 'en',
  interpolation: { escapeValue: false },
});

i18n.on('languageChanged', (lng) => {
  if (typeof window !== 'undefined') {
    localStorage.setItem('locale', lng);
  }
});

export function getSavedLocale(): string | null {
  if (typeof window === 'undefined') return null;
  const saved = localStorage.getItem('locale');
  if (saved) return saved;
  const browser = navigator.language.split('-')[0];
  return browser === 'ru' ? 'ru' : 'en';
}

export default i18n;
