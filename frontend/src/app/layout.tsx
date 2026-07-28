import type { Metadata } from 'next';
import { Inter } from 'next/font/google';
import { ThemeProvider } from 'next-themes';
import { QueryClientProvider } from '@/components/providers/query-client-provider';
import { AuthProvider } from '@/hooks/use-auth';
import { I18nProvider } from '@/components/providers/i18n-provider';
import '@/styles/globals.css';

const inter = Inter({ subsets: ['latin'], variable: '--font-inter' });

export const metadata: Metadata = {
  title: 'DevFlow',
  description: 'Project management platform for development teams',
  icons: { icon: '/favicon.ico' },
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className={`${inter.variable} font-sans antialiased`}>
        <QueryClientProvider>
          <AuthProvider>
            <I18nProvider>
              <ThemeProvider
                attribute="class"
                defaultTheme="system"
                enableSystem
                disableTransitionOnChange
              >
                {children}
              </ThemeProvider>
            </I18nProvider>
          </AuthProvider>
        </QueryClientProvider>
      </body>
    </html>
  );
}
