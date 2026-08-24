import type { ForgeConfig } from '@electron-forge/shared-types';
import { MakerSquirrel } from '@electron-forge/maker-squirrel';
import { MakerZIP } from '@electron-forge/maker-zip';
import { MakerDMG } from '@electron-forge/maker-dmg';
import { MakerDeb } from '@electron-forge/maker-deb';
import { MakerRpm } from '@electron-forge/maker-rpm';
import { VitePlugin } from '@electron-forge/plugin-vite';
import { FusesPlugin } from '@electron-forge/plugin-fuses';
import { FuseV1Options, FuseVersion } from '@electron/fuses';

const config: ForgeConfig = {
  packagerConfig: {
    asar: true,
    name: 'Look News',
    executableName: 'look-news',
    icon: './assets/icon',
    appBundleId: 'com.clemilsonazevedo.looknews',
    appCategoryType: 'public.app-category.productivity',
    extraResource: ['./assets'],
  },
  rebuildConfig: {},
  makers: [
    new MakerSquirrel({
      name: 'Look News',
      authors: 'Clemilson Azevedo',
      description: 'Agregador de noticias na sua barra de menu',
      setupExe: 'look-news.exe',
      setupIcon: './assets/trayIconTemplate.png',
    }),

    new MakerDMG({
      name: 'Look News',
      format: 'ULFO',
    }),
    new MakerZIP({}, ['darwin']),

    new MakerDeb({
      options: {
        name: 'look-news',
        productName: 'Look News',
        genericName: 'Look News',
        description: 'Agregador de noticias na sua barra de menu',
        categories: ['Utility'],
        maintainer: 'Clemilson Azevedo <clemilsondeazevedo@gmail.com>',
        homepage: 'https://github.com/clemilsonazevedo/look-news',
      },
    }),
    new MakerRpm({
      options: {
        name: 'look-news',
        productName: 'Look News',
        description: 'Agregador de noticias na sua barra de menu',
        categories: ['Utility'],
        homepage: 'https://github.com/clemilsonazevedo/look-news',
      },
    }),
  ],
  plugins: [
    new VitePlugin({
      build: [
        {
          entry: 'src/main.ts',
          config: 'vite.main.config.ts',
          target: 'main',
        },
        {
          entry: 'src/preload.ts',
          config: 'vite.preload.config.ts',
          target: 'preload',
        },
      ],
      renderer: [
        {
          name: 'main_window',
          config: 'vite.renderer.config.ts',
        },
      ],
    }),
    new FusesPlugin({
      version: FuseVersion.V1,
      [FuseV1Options.RunAsNode]: false,
      [FuseV1Options.EnableCookieEncryption]: true,
      [FuseV1Options.EnableNodeOptionsEnvironmentVariable]: false,
      [FuseV1Options.EnableNodeCliInspectArguments]: false,
      [FuseV1Options.EnableEmbeddedAsarIntegrityValidation]: true,
      [FuseV1Options.OnlyLoadAppFromAsar]: true,
    }),
  ],
};

export default config;