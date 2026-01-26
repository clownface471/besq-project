export interface Chart {
  id: string;
  x: number;
  y: number;
  width: number;
  height: number;
  title: string;
  color?: string;
}

export interface Arrow {
  id: string;
  fromId: string;
  toId: string;
  curvature: number;
  color?: string;
}