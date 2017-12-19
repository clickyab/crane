export default interface IAd {
  element: Element,
  slot: string | null,
  width: string | null,
  height: string | null,
  adType?: string | null,
  minFlex?: string | null,
  size?: number;
  valid?: boolean;
  src?: string;
  effect?: string | null;
  iframe?: string;
}
