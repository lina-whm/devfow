import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import en from './en.json';
import ru from './ru.json';

const resources = { en: { translation: en }, ru: { translation: ru } };

const savedLang = typeof window !== 'undefined' ? localStorage.getItem('locale') : null;
const browserLang = typeof window !== 'undefined' ? navigator.language.split('-')[0] : 'en';
const defaultLang = savedLang || (browserLang === 'ru' ? 'ru' : 'en');

i18n.use(initReactI18next).init({
  resources,
  lng: defaultLang,
  fallbackLng: 'en',
  interpolation: { escapeValue: false },
});

export default i18n;
