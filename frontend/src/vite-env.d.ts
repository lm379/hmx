/// <reference types="vite/client" />

declare module '*.svg?component' {
  import { FunctionalComponent, SVGAttributes } from 'vue'
  const src: FunctionalComponent<SVGAttributes>
  export default src
}

declare module 'pdfjs-dist/build/pdf.mjs' {
  const pdfjs: any
  export = pdfjs
}
