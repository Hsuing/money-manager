import type { CapacitorConfig } from '@capacitor/cli';

const config: CapacitorConfig = {
  appId: 'com.hsuing.accounts',
  appName: 'EasyAccounts',
  webDir: 'dist',
  plugins: {
    SplashScreen: {
      launchShowDuration: 2000,
      launchAutoHide: true,
      backgroundColor: "#080808",
      androidSplashResourceName: "splash",
      androidScaleType: "CENTER",
      showSpinner: false,
    },
  },
};

export default config;
