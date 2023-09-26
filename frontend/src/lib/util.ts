import {
  type SidebarItem,
  TabIndex,
  TabContent,
  type KubeGVRK,
  type HTTPStatusCode,
} from "./types";

export const httpStatus: HTTPStatusCode = {
  200: "OK",
  201: "Created",
  400: "Bad Request",
  401: "Unauthorized",
  403: "Forbidden",
  404: "Not Found",
  405: "Method Not Allowed",
  408: "Request Timeout",
  500: "Internal Server Error",
  501: "Not Implemented",
  502: "Bad Gateway",
  503: "Service Unavailable",
  504: "Gateway Timeout",
};

export const apiURL = {
  auth: "/auth",
  data: "/srv/data",
  dataWS: "/srv/data/ws",
  logs: "/srv/logs",
  logsStream: "/srv/logs/stream",
  shell: "/srv/shell",
  shellExec: "/srv/shell/exec",
} as const;

export const colorScheme = {
  light: {
    primary: "#65cc89",
    "primary-focus": "#4bc374",
    "primary-content": "#233e31",
    secondary: "#367bfb",
    "secondary-focus": "#1364fa",
    "secondary-content": "#f9fafb",
    accent: "#e95234",
    "accent-focus": "#e13918",
    "accent-content": "#f9fafb",
    neutral: "#333c4c",
    "neutral-focus": "#252b37",
    "neutral-content": "#f9fafb",
    "base-100": "#e9e7e7",
    "base-200": "#d7d5d5",
    "base-300": "#c7c2c2",
    "base-content": "#333c4c",
    info: "#3abef7",
    "info-content": "#012b3e",
    success: "#37d39a",
    "success-content": "#013321",
    warning: "#fabd22",
    "warning-content": "#382800",
    error: "#f87272",
    "error-content": "#470000",
  },
  dark: {
    primary: "#065f46",
    "primary-focus": "#199945",
    "primary-content": "#ffffff",
    secondary: "#0ea5e9",
    "secondary-focus": "#38bdf8",
    "secondary-content": "#bae6fd",
    accent: "#fb7185",
    "accent-focus": "#e11d48",
    "accent-content": "#ffffff",
    neutral: "#414557",
    "neutral-focus": "#313543",
    "neutral-content": "#d6d7db",
    "base-100": "#262935",
    "base-200": "#181921",
    "base-300": "#09090d",
    "base-content": "#f8f8f1",
    info: "#8be8fd",
    "info-content": "#202e31",
    success: "#52f97c",
    "success-content": "#192e1c",
    // warning: '#f0fa89',
    warning: "#db7706",
    "warning-content": "#2d2e1e",
    error: "#ff5756",
    "error-content": "#311816",
  },
} as const;

export const icons = {
  settings:
    "M10.5 6h9.75M10.5 6a1.5 1.5 0 11-3 0m3 0a1.5 1.5 0 10-3 0M3.75 6H7.5m3 12h9.75m-9.75 0a1.5 1.5 0 01-3 0m3 0a1.5 1.5 0 00-3 0m-3.75 0H7.5m9-6h3.75m-3.75 0a1.5 1.5 0 01-3 0m3 0a1.5 1.5 0 00-3 0m-9.75 0h9.75",
  close: "M6 18L18 6M6 6l12 12",
  closeCircle:
    "M9.75 9.75l4.5 4.5m0-4.5l-4.5 4.5M21 12a9 9 0 11-18 0 9 9 0 0118 0z",
  cogWheel:
    "M9.594 3.94c.09-.542.56-.94 1.11-.94h2.593c.55 0 1.02.398 1.11.94l.213 1.281c.063.374.313.686.645.87.074.04.147.083.22.127.324.196.72.257 1.075.124l1.217-.456a1.125 1.125 0 011.37.49l1.296 2.247a1.125 1.125 0 01-.26 1.431l-1.003.827c-.293.24-.438.613-.431.992a6.759 6.759 0 010 .255c-.007.378.138.75.43.99l1.005.828c.424.35.534.954.26 1.43l-1.298 2.247a1.125 1.125 0 01-1.369.491l-1.217-.456c-.355-.133-.75-.072-1.076.124a6.57 6.57 0 01-.22.128c-.331.183-.581.495-.644.869l-.213 1.28c-.09.543-.56.941-1.11.941h-2.594c-.55 0-1.02-.398-1.11-.94l-.213-1.281c-.062-.374-.312-.686-.644-.87a6.52 6.52 0 01-.22-.127c-.325-.196-.72-.257-1.076-.124l-1.217.456a1.125 1.125 0 01-1.369-.49l-1.297-2.247a1.125 1.125 0 01.26-1.431l1.004-.827c.292-.24.437-.613.43-.992a6.932 6.932 0 010-.255c.007-.378-.138-.75-.43-.99l-1.004-.828a1.125 1.125 0 01-.26-1.43l1.297-2.247a1.125 1.125 0 011.37-.491l1.216.456c.356.133.751.072 1.076-.124.072-.044.146-.087.22-.128.332-.183.582-.495.644-.869l.214-1.281z M15 12a3 3 0 11-6 0 3 3 0 016 0z",
  user: "M17.982 18.725A7.488 7.488 0 0012 15.75a7.488 7.488 0 00-5.982 2.975m11.963 0a9 9 0 10-11.963 0m11.963 0A8.966 8.966 0 0112 21a8.966 8.966 0 01-5.982-2.275M15 9.75a3 3 0 11-6 0 3 3 0 016 0z",
  // sun: 'M5.64,17l-.71.71a1,1,0,0,0,0,1.41,1,1,0,0,0,1.41,0l.71-.71A1,1,0,0,0,5.64,17ZM5,12a1,1,0,0,0-1-1H3a1,1,0,0,0,0,2H4A1,1,0,0,0,5,12Zm7-7a1,1,0,0,0,1-1V3a1,1,0,0,0-2,0V4A1,1,0,0,0,12,5ZM5.64,7.05a1,1,0,0,0,.7.29,1,1,0,0,0,.71-.29,1,1,0,0,0,0-1.41l-.71-.71A1,1,0,0,0,4.93,6.34Zm12,.29a1,1,0,0,0,.7-.29l.71-.71a1,1,0,1,0-1.41-1.41L17,5.64a1,1,0,0,0,0,1.41A1,1,0,0,0,17.66,7.34ZM21,11H20a1,1,0,0,0,0,2h1a1,1,0,0,0,0-2Zm-9,8a1,1,0,0,0-1,1v1a1,1,0,0,0,2,0V20A1,1,0,0,0,12,19ZM18.36,17A1,1,0,0,0,17,18.36l.71.71a1,1,0,0,0,1.41,0,1,1,0,0,0,0-1.41ZM12,6.5A5.5,5.5,0,1,0,17.5,12,5.51,5.51,0,0,0,12,6.5Zm0,9A3.5,3.5,0,1,1,15.5,12,3.5,3.5,0,0,1,12,15.5Z',
  sun: "M12 3v2.25m6.364.386l-1.591 1.591M21 12h-2.25m-.386 6.364l-1.591-1.591M12 18.75V21m-4.773-4.227l-1.591 1.591M5.25 12H3m4.227-4.773L5.636 5.636M15.75 12a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0z",
  // moon: 'M21.64,13a1,1,0,0,0-1.05-.14,8.05,8.05,0,0,1-3.37.73A8.15,8.15,0,0,1,9.08,5.49a8.59,8.59,0,0,1,.25-2A1,1,0,0,0,8,2.36,10.14,10.14,0,1,0,22,14.05,1,1,0,0,0,21.64,13Zm-9.5,6.69A8.14,8.14,0,0,1,7.08,5.22v.27A10.15,10.15,0,0,0,17.22,15.63a9.79,9.79,0,0,0,2.1-.22A8.11,8.11,0,0,1,12.14,19.73Z',
  moon: "M21.752 15.002A9.718 9.718 0 0118 15.75c-5.385 0-9.75-4.365-9.75-9.75 0-1.33.266-2.597.748-3.752A9.753 9.753 0 003 11.25C3 16.635 7.365 21 12.75 21a9.753 9.753 0 009.002-5.998z",
  notificationBell:
    "M14.857 17.082a23.848 23.848 0 005.454-1.31A8.967 8.967 0 0118 9.75v-.7V9A6 6 0 006 9v.75a8.967 8.967 0 01-2.312 6.022c1.733.64 3.56 1.085 5.455 1.31m5.714 0a24.255 24.255 0 01-5.714 0m5.714 0a3 3 0 11-5.714 0",
  sidebarMenuButton:
    "M3.75 6A2.25 2.25 0 016 3.75h2.25A2.25 2.25 0 0110.5 6v2.25a2.25 2.25 0 01-2.25 2.25H6a2.25 2.25 0 01-2.25-2.25V6zM3.75 15.75A2.25 2.25 0 016 13.5h2.25a2.25 2.25 0 012.25 2.25V18a2.25 2.25 0 01-2.25 2.25H6A2.25 2.25 0 013.75 18v-2.25zM13.5 6a2.25 2.25 0 012.25-2.25H18A2.25 2.25 0 0120.25 6v2.25A2.25 2.25 0 0118 10.5h-2.25a2.25 2.25 0 01-2.25-2.25V6zM13.5 15.75a2.25 2.25 0 012.25-2.25H18a2.25 2.25 0 012.25 2.25V18A2.25 2.25 0 0118 20.25h-2.25A2.25 2.25 0 0113.5 18v-2.25z",
  home: "M2.25 12l8.954-8.955c.44-.439 1.152-.439 1.591 0L21.75 12M4.5 9.75v10.125c0 .621.504 1.125 1.125 1.125H9.75v-4.875c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125V21h4.125c.621 0 1.125-.504 1.125-1.125V9.75M8.25 21h8.25",
  table: "M3.75 5.25h16.5m-16.5 4.5h16.5m-16.5 4.5h16.5m-16.5 4.5h16.5",
  graph:
    "M13.5 10.5H21A7.5 7.5 0 0013.5 3v7.5zM10.5 6a7.5 7.5 0 107.5 7.5h-7.5V6z",
  nodes:
    "M8.25 3v1.5M4.5 8.25H3m18 0h-1.5M4.5 12H3m18 0h-1.5m-15 3.75H3m18 0h-1.5M8.25 19.5V21M12 3v1.5m0 15V21m3.75-18v1.5m0 15V21m-9-1.5h10.5a2.25 2.25 0 002.25-2.25V6.75a2.25 2.25 0 00-2.25-2.25H6.75A2.25 2.25 0 004.5 6.75v10.5a2.25 2.25 0 002.25 2.25zm.75-12h9v9h-9v-9z",
  // 'M17 7H15V9H17V7Z M7 11H9V13H7V11Z M13 11H11V13H13V11Z M15 11H17V13H15V11Z M9 15H7V17H9V15Z M11 15H13V17H11V15Z M17 15H15V17H17V15Z',
  about:
    "M7.5 8.25h9m-9 3H12m-9.75 1.51c0 1.6 1.123 2.994 2.707 3.227 1.129.166 2.27.293 3.423.379.35.026.67.21.865.501L12 21l2.755-4.133a1.14 1.14 0 01.865-.501 48.172 48.172 0 003.423-.379c1.584-.233 2.707-1.626 2.707-3.228V6.741c0-1.602-1.123-2.995-2.707-3.228A48.394 48.394 0 0012 3c-2.392 0-4.744.175-7.043.513C3.373 3.746 2.25 5.14 2.25 6.741v6.018z",
  notFound:
    "M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636",
  pods: "M21 7.5l-9-5.25L3 7.5m18 0l-9 5.25m9-5.25v9l-9 5.25M3 7.5l9 5.25M3 7.5v9l9 5.25m0-9v9",
  deployment:
    "M15.59 14.37a6 6 0 01-5.84 7.38v-4.8m5.84-2.58a14.98 14.98 0 006.16-12.12A14.98 14.98 0 009.631 8.41m5.96 5.96a14.926 14.926 0 01-5.841 2.58m-.119-8.54a6 6 0 00-7.381 5.84h4.8m2.581-5.84a14.927 14.927 0 00-2.58 5.84m2.699 2.7c-.103.021-.207.041-.311.06a15.09 15.09 0 01-2.448-2.448 14.9 14.9 0 01.06-.312m-2.24 2.39a4.493 4.493 0 00-1.757 4.306 4.493 4.493 0 004.306-1.758M16.5 9a1.5 1.5 0 11-3 0 1.5 1.5 0 013 0z",
  //   arrowLeft: 'M9 15L3 9m0 0l6-6M3 9h12a6 6 0 010 12h-3',
  //   arrowRight: 'M15 15l6-6m0 0l-6-6m6 6H9a6 6 0 000 12h3'
  arrowLeft: "M19.5 12h-15m0 0l6.75 6.75M4.5 12l6.75-6.75",
  arrowRight: "M4.5 12h15m0 0l-6.75-6.75M19.5 12l-6.75 6.75",
  helm: "M12 9v12m-8 -8a8 8 0 0 0 16 0m1 0h-2m-14 0h-2 M12 6m-3 0a3 3 0 1 0 6 0a3 3 0 1 0 -6 0",
  document:
    "M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m2.25 0H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z",
  documentList:
    "M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z",
  documentDetails:
    "M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m5.231 13.481L15 17.25m-4.5-15H5.625c-.621 0-1.125.504-1.125 1.125v16.5c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9zm3.75 11.625a2.625 2.625 0 11-5.25 0 2.625 2.625 0 015.25 0z",
  code: "M17.25 6.75L22.5 12l-5.25 5.25m-10.5 0L1.5 12l5.25-5.25m7.5-3l-4.5 16.5",
  codeSquared:
    "M14.25 9.75L16.5 12l-2.25 2.25m-4.5 0L7.5 12l2.25-2.25M6 20.25h12A2.25 2.25 0 0020.25 18V6A2.25 2.25 0 0018 3.75H6A2.25 2.25 0 003.75 6v12A2.25 2.25 0 006 20.25z",
  cmdLine:
    "M8 9l3 3l-3 3 M13 15l3 0 M3 4m0 2a2 2 0 0 1 2 -2h14a2 2 0 0 1 2 2v12a2 2 0 0 1 -2 2h-14a2 2 0 0 1 -2 -2z",
  terminal:
    "M9.26,10.23c.98,.98,.98,2.57,0,3.54L.85,21.86c-.1,.09-.22,.14-.35,.14-.13,0-.26-.05-.36-.15-.19-.2-.19-.52,.01-.71L8.57,13.05c.58-.58,.58-1.53,0-2.11L.15,2.86c-.2-.19-.21-.51-.01-.71,.19-.2,.51-.21,.71-.01L9.26,10.23Zm14.24,10.77H10.5c-.28,0-.5,.22-.5,.5s.22,.5,.5,.5h13c.28,0,.5-.22,.5-.5s-.22-.5-.5-.5Z",
  event:
    "M7.5 8.25h9m-9 3H12m-9.75 1.51c0 1.6 1.123 2.994 2.707 3.227 1.129.166 2.27.293 3.423.379.35.026.67.21.865.501L12 21l2.755-4.133a1.14 1.14 0 01.865-.501 48.172 48.172 0 003.423-.379c1.584-.233 2.707-1.626 2.707-3.228V6.741c0-1.602-1.123-2.995-2.707-3.228A48.394 48.394 0 0012 3c-2.392 0-4.744.175-7.043.513C3.373 3.746 2.25 5.14 2.25 6.741v6.018z",
  queue:
    "M3.75 12h16.5m-16.5 3.75h16.5M3.75 19.5h16.5M5.625 4.5h12.75a1.875 1.875 0 010 3.75H5.625a1.875 1.875 0 010-3.75z",
  stack:
    "M6.429 9.75L2.25 12l4.179 2.25m0-4.5l5.571 3 5.571-3m-11.142 0L2.25 7.5 12 2.25l9.75 5.25-4.179 2.25m0 0L21.75 12l-4.179 2.25m0 0l4.179 2.25L12 21.75 2.25 16.5l4.179-2.25m11.142 0l-5.571 3-5.571-3",
  newspaper:
    "M12 7.5h1.5m-1.5 3h1.5m-7.5 3h7.5m-7.5 3h7.5m3-9h3.375c.621 0 1.125.504 1.125 1.125V18a2.25 2.25 0 01-2.25 2.25M16.5 7.5V18a2.25 2.25 0 002.25 2.25M16.5 7.5V4.875c0-.621-.504-1.125-1.125-1.125H4.125C3.504 3.75 3 4.254 3 4.875V18a2.25 2.25 0 002.25 2.25h13.5M6 7.5h3v3H6v-3z",
  info: "M11.25 11.25l.041-.02a.75.75 0 011.063.852l-.708 2.836a.75.75 0 001.063.853l.041-.021M21 12a9 9 0 11-18 0 9 9 0 0118 0zm-9-3.75h.008v.008H12V8.25z",
  warning:
    "M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126zM12 15.75h.007v.008H12v-.008z",
  clock: "M12 6v6h4.5m4.5 0a9 9 0 11-18 0 9 9 0 0118 0z",
  elipsisVertical:
    "M12 6.75a.75.75 0 110-1.5.75.75 0 010 1.5zM12 12.75a.75.75 0 110-1.5.75.75 0 010 1.5zM12 18.75a.75.75 0 110-1.5.75.75 0 010 1.5z",
  ingress:
    "M7.5 21L3 16.5m0 0L7.5 12M3 16.5h13.5m0-13.5L21 7.5m0 0L16.5 12M21 7.5H7.5",
  statefulSet:
    "M20.25 6.375c0 2.278-3.694 4.125-8.25 4.125S3.75 8.653 3.75 6.375m16.5 0c0-2.278-3.694-4.125-8.25-4.125S3.75 4.097 3.75 6.375m16.5 0v11.25c0 2.278-3.694 4.125-8.25 4.125s-8.25-1.847-8.25-4.125V6.375m16.5 0v3.75m-16.5-3.75v3.75m16.5 0v3.75C20.25 16.153 16.556 18 12 18s-8.25-1.847-8.25-4.125v-3.75m16.5 0c0 2.278-3.694 4.125-8.25 4.125s-8.25-1.847-8.25-4.125",
  // job: 'M3 17H21V19H3V17ZM3 11H6V14H3V11ZM8 11H11V14H8V11ZM3 5H6V8H3V5ZM13 5H16V8H13V5ZM18 5H21V8H18V5ZM13 11H16V14H13V11ZM18 11H21V14H18V11ZM8 5H11V8H8V5Z'
  // job: 'M3 17H21V19H3V17ZM3 11H6V1 11H11V14H8 5H6V8H3V5ZM13 5H16V8H13V5ZM18 5H21V8H18V5ZM13 11H16V14H13V11ZM18 11H21V14H18V11ZM8 5H11V8H8V5Z'
  job: "M5.5 8H8.5V11H5.5V8Z M5.5 13V16H8.5V13H5.5Z M10.5 8V11H13.5V8H10.5Z M10.5 13H13.5V16H10.5V13Z M15.5 3H18.5V6H15.5V3Z M15.5 8V11H18.5V8H15.5Z M15.5 13H18.5V16H15.5V13Z",
  tag: "M9.568 3H5.25A2.25 2.25 0 003 5.25v4.318c0 .597.237 1.17.659 1.591l9.581 9.581c.699.699 1.78.872 2.607.33a18.095 18.095 0 005.223-5.223c.542-.827.369-1.908-.33-2.607L11.16 3.66A2.25 2.25 0 009.568 3z M6 6h.008v.008H6V6z",
  hashtag: "M5.25 8.25h15m-16.5 7.5h15m-1.8-13.5l-3.9 19.5m-2.1-19.5l-3.9 19.5",
  bookmark:
    "M16.5 3.75V16.5L12 14.25 7.5 16.5V3.75m9 0H18A2.25 2.25 0 0120.25 6v12A2.25 2.25 0 0118 20.25H6A2.25 2.25 0 013.75 18V6A2.25 2.25 0 016 3.75h1.5m9 0h-9",
  identification:
    "M15 9h3.75M15 12h3.75M15 15h3.75M4.5 19.5h15a2.25 2.25 0 002.25-2.25V6.75A2.25 2.25 0 0019.5 4.5h-15a2.25 2.25 0 00-2.25 2.25v10.5A2.25 2.25 0 004.5 19.5zm6-10.125a1.875 1.875 0 11-3.75 0 1.875 1.875 0 013.75 0zm1.294 6.336a6.721 6.721 0 01-3.17.789 6.721 6.721 0 01-3.168-.789 3.376 3.376 0 016.338 0z",
  paperAirplane:
    "M6 12L3.269 3.126A59.768 59.768 0 0121.485 12 59.77 59.77 0 013.27 20.876L5.999 12zm0 0h7.5",
  secret:
    "M16.5 10.5V6.75a4.5 4.5 0 10-9 0v3.75m-.75 11.25h10.5a2.25 2.25 0 002.25-2.25v-6.75a2.25 2.25 0 00-2.25-2.25H6.75a2.25 2.25 0 00-2.25 2.25v6.75a2.25 2.25 0 002.25 2.25z",
  configMap:
    "M9 6.75V15m6-6v8.25m.503 3.498l4.875-2.437c.381-.19.622-.58.622-1.006V4.82c0-.836-.88-1.38-1.628-1.006l-3.869 1.934c-.317.159-.69.159-1.006 0L9.503 3.252a1.125 1.125 0 00-1.006 0L3.622 5.689C3.24 5.88 3 6.27 3 6.695V19.18c0 .836.88 1.38 1.628 1.006l3.869-1.934c.317-.159.69-.159 1.006 0l4.994 2.497c.317.158.69.158 1.006 0z",
  at: "M16.5 12a4.5 4.5 0 11-9 0 4.5 4.5 0 019 0zm0 0c0 1.657 1.007 3 2.25 3S21 13.657 21 12a9 9 0 10-2.636 6.364M16.5 12V8.25",
  trash:
    "M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0",
  listBullet:
    "M8.25 6.75h12M8.25 12h12m-12 5.25h12M3.75 6.75h.007v.008H3.75V6.75zm.375 0a.375.375 0 11-.75 0 .375.375 0 01.75 0zM3.75 12h.007v.008H3.75V12zm.375 0a.375.375 0 11-.75 0 .375.375 0 01.75 0zm-.375 5.25h.007v.008H3.75v-.008zm.375 0a.375.375 0 11-.75 0 .375.375 0 01.75 0z",
  arrowDown: "M19.5 5.25l-7.5 7.5-7.5-7.5m15 6l-7.5 7.5-7.5-7.5",
  search:
    "M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z",
  dashboard:
    "M3.75 6A2.25 2.25 0 016 3.75h2.25A2.25 2.25 0 0110.5 6v2.25a2.25 2.25 0 01-2.25 2.25H6a2.25 2.25 0 01-2.25-2.25V6zM3.75 15.75A2.25 2.25 0 016 13.5h2.25a2.25 2.25 0 012.25 2.25V18a2.25 2.25 0 01-2.25 2.25H6A2.25 2.25 0 013.75 18v-2.25zM13.5 6a2.25 2.25 0 012.25-2.25H18A2.25 2.25 0 0120.25 6v2.25A2.25 2.25 0 0118 10.5h-2.25a2.25 2.25 0 01-2.25-2.25V6zM13.5 15.75a2.25 2.25 0 012.25-2.25H18a2.25 2.25 0 012.25 2.25V18A2.25 2.25 0 0118 20.25h-2.25A2.25 2.25 0 0113.5 18v-2.25z",
  filter:
    "M10.5 6h9.75M10.5 6a1.5 1.5 0 11-3 0m3 0a1.5 1.5 0 10-3 0M3.75 6H7.5m3 12h9.75m-9.75 0a1.5 1.5 0 01-3 0m3 0a1.5 1.5 0 00-3 0m-3.75 0H7.5m9-6h3.75m-3.75 0a1.5 1.5 0 01-3 0m3 0a1.5 1.5 0 00-3 0m-9.75 0h9.75",
} as const;

export const tabs: TabContent = {
  details: {
    index: TabIndex.DETAILS,
    name: "Details",
    icon: icons.documentDetails,
    strokeWidth: 1,
    viewBox: "0 0 24 24",
  },
  yaml: {
    index: TabIndex.YAML,
    name: "YAML",
    icon: icons.code,
    strokeWidth: 1,
    viewBox: "0 0 24 24",
  },
  logs: {
    index: TabIndex.LOGS,
    name: "Logs",
    icon: icons.documentList,
    strokeWidth: 1,
    viewBox: "0 0 24 24",
  },
  shell: {
    index: TabIndex.SHELL,
    name: "Shell",
    icon: icons.cmdLine,
    strokeWidth: 1,
    viewBox: "0 0 24 24",
  },
  events: {
    index: TabIndex.EVENTS,
    name: "Events",
    icon: icons.newspaper,
    strokeWidth: 1,
    viewBox: "0 0 24 24",
  },
  queue: {
    index: TabIndex.QUEUE,
    name: "Queue",
    icon: icons.queue,
    strokeWidth: 1,
    viewBox: "0 0 24 24",
  },
  stack: {
    index: TabIndex.STACK,
    name: "Stack",
    icon: icons.stack,
    strokeWidth: 1,
    viewBox: "0 0 24 24",
  },
  search: {
    index: TabIndex.SEARCH,
    name: "Search",
    icon: icons.search,
    strokeWidth: 1,
    viewBox: "0 0 24 24",
  },
  dashboard: {
    index: TabIndex.DASHBOARD,
    name: "Dashboard",
    icon: icons.dashboard,
    strokeWidth: 1,
    viewBox: "0 0 24 24",
  },
} as const;

export const routeString = {
  home: "/",
  nodeList: "/node",
  podList: "/pod",
  deployList: "/deployment",
  serviceList: "/service",
  secretList: "/secret",
  configMapList: "/configmap",
  statefulSetList: "/statefulset",
  ingressList: "/ingress",
  jobList: "/job",
  about: "/about",
  hello: "/hello/svelte",
  wildcard: "/wild/card",
  notExist: "/does/not/exist",
  helm: "/release",
  clusterRoleList: "/clusterrole",
  clusterRoleBindingList: "/clusterrolebinding",
  roleList: "/role",
  roleBindingList: "/rolebinding",
  cronJobList: "/cronjob",
  daemonSetList: "/daemonset",
  eventList: "/event",
  ingressClassList: "/ingressclass",
  namespaceList: "/namespace",
  networkPolicyList: "/networkpolicy",
  persistentVolumeList: "/persistentvolume",
  persistentVolumeClaimList: "/persistentvolumeclaim",
  replicationControllerList: "/replicationcontroller",
  serviceAccountList: "/serviceaccount",
  storageClassList: "/storageclass",
} as const;

export const sidebarItems: SidebarItem[] = [
  {
    name: "Home",
    icon: icons.home,
    url: routeString.home,
    viewBox: "0 0 24 24",
    strokeWidth: 1.25,
  },
  {
    name: "Helm Releases",
    icon: icons.helm,
    url: routeString.helm,
    viewBox: "0 0 24 24",
    strokeWidth: 1.25,
  },
  {
    name: "Nodes",
    icon: icons.nodes,
    url: routeString.nodeList,
    viewBox: "0 0 24 24",
    strokeWidth: 1.25,
  },
  {
    name: "Pods",
    icon: icons.pods,
    url: routeString.podList,
    viewBox: "0 0 24 24",
    strokeWidth: 1.25,
  },
  {
    name: "Ingresses",
    icon: icons.ingress,
    url: routeString.ingressList,
    viewBox: "0 0 24 24",
    strokeWidth: 1.25,
  },
  {
    name: "Deployments",
    icon: icons.deployment,
    url: routeString.deployList,
    viewBox: "0 0 24 24",
    strokeWidth: 1.25,
  },
  {
    name: "Services",
    icon: icons.hashtag,
    url: routeString.serviceList,
    viewBox: "0 0 24 24",
    strokeWidth: 1.25,
  },
  {
    name: "Secrets",
    icon: icons.secret,
    url: routeString.secretList,
    viewBox: "0 0 24 24",
    strokeWidth: 1.25,
  },
  {
    name: "ConfigMaps",
    icon: icons.configMap,
    url: routeString.configMapList,
    viewBox: "0 0 24 24",
    strokeWidth: 1.25,
  },
  {
    name: "StatefulSets",
    icon: icons.queue,
    url: routeString.statefulSetList,
    viewBox: "0 0 24 24",
    strokeWidth: 1.25,
  },
  {
    name: "Jobs",
    icon: icons.job,
    url: routeString.jobList,
    strokeWidth: 0.8,
    viewBox: "0 0 22 22",
  },
  {
    name: "About",
    icon: icons.about,
    url: routeString.about,
    viewBox: "0 0 24 24",
    strokeWidth: 1.25,
  },
  {
    name: "Doesn't exist",
    icon: icons.notFound,
    url: routeString.notExist,
    viewBox: "0 0 24 24",
    strokeWidth: 1.25,
  },
  {
    name: "ClusterRoles",
    icon: icons.hashtag,
    url: routeString.clusterRoleList,
    viewBox: "0 0 24 24",
    strokeWidth: 1.25,
  },
  {
    name: "ClusterRoleBindings",
    icon: icons.hashtag,
    url: routeString.clusterRoleBindingList,
    viewBox: "0 0 24 24",
    strokeWidth: 1.25,
  },
  {
    name: "Roles",
    icon: icons.hashtag,
    url: routeString.roleList,
    viewBox: "0 0 24 24",
    strokeWidth: 1.25,
  },
  {
    name: "RoleBindings",
    icon: icons.hashtag,
    url: routeString.roleBindingList,
    viewBox: "0 0 24 24",
    strokeWidth: 1.25,
  },
  {
    name: "CronJobs",
    icon: icons.hashtag,
    url: routeString.cronJobList,
    viewBox: "0 0 24 24",
    strokeWidth: 1.25,
  },
  {
    name: "DaemonSets",
    icon: icons.hashtag,
    url: routeString.daemonSetList,
    viewBox: "0 0 24 24",
    strokeWidth: 1.25,
  },
  {
    name: "Events",
    icon: icons.hashtag,
    url: routeString.eventList,
    viewBox: "0 0 24 24",
    strokeWidth: 1.25,
  },
  {
    name: "IngressClasses",
    icon: icons.hashtag,
    url: routeString.ingressClassList,
    viewBox: "0 0 24 24",
    strokeWidth: 1.25,
  },
  {
    name: "Namespaces",
    icon: icons.hashtag,
    url: routeString.namespaceList,
    viewBox: "0 0 24 24",
    strokeWidth: 1.25,
  },
  {
    name: "NetworkPolicies",
    icon: icons.hashtag,
    url: routeString.networkPolicyList,
    viewBox: "0 0 24 24",
    strokeWidth: 1.25,
  },
  {
    name: "PersistentVolumes",
    icon: icons.hashtag,
    url: routeString.persistentVolumeList,
    viewBox: "0 0 24 24",
    strokeWidth: 1.25,
  },
  {
    name: "PersistentVolumeClaims",
    icon: icons.hashtag,
    url: routeString.persistentVolumeClaimList,
    viewBox: "0 0 24 24",
    strokeWidth: 1.25,
  },
  {
    name: "ReplicationController",
    icon: icons.hashtag,
    url: routeString.replicationControllerList,
    viewBox: "0 0 24 24",
    strokeWidth: 1.25,
  },
  {
    name: "ServiceAccounts",
    icon: icons.hashtag,
    url: routeString.serviceAccountList,
    viewBox: "0 0 24 24",
    strokeWidth: 1.25,
  },
  {
    name: "StorageClasses",
    icon: icons.hashtag,
    url: routeString.storageClassList,
    viewBox: "0 0 24 24",
    strokeWidth: 1.25,
  },
];

export const KubeResourceType: Array<KubeGVRK> = [
  {
    group: "",
    version: "v1",
    resource: "namespaces",
    kind: "Namespace",
    isNamespaced: false,
  },

  {
    group: "",
    version: "v1",
    resource: "nodes",
    kind: "Node",
    isNamespaced: false,
  },

  {
    group: "metrics.k8s.io",
    version: "v1beta1",
    resource: "nodes",
    kind: "NodeMetrics",
    isNamespaced: false,
  },

  {
    group: "",
    version: "v1",
    resource: "persistentvolumes",
    kind: "PersistentVolume",
    isNamespaced: false,
  },

  {
    group: "",
    version: "v1",
    resource: "persistentvolumeclaims",
    kind: "PersistentVolumeClaim",
    isNamespaced: true,
  },

  {
    group: "",
    version: "v1",
    resource: "replicationcontrollers",
    kind: "ReplicationController",
    isNamespaced: true,
  },

  {
    group: "",
    version: "v1",
    resource: "pods",
    kind: "Pod",
    isNamespaced: true,
  },

  {
    group: "metrics.k8s.io",
    version: "v1beta1",
    resource: "pods",
    kind: "PodMetrics",
    isNamespaced: true,
  },

  {
    group: "",
    version: "v1",
    resource: "services",
    kind: "Service",
    isNamespaced: true,
  },

  {
    group: "",
    version: "v1",
    resource: "secrets",
    kind: "Secret",
    isNamespaced: true,
  },

  {
    group: "",
    version: "v1",
    resource: "serviceaccounts",
    kind: "ServiceAccount",
    isNamespaced: true,
  },

  {
    group: "",
    version: "v1",
    resource: "configmaps",
    kind: "ConfigMap",
    isNamespaced: true,
  },

  {
    group: "apps",
    version: "v1",
    resource: "deployments",
    kind: "Deployment",
    isNamespaced: true,
  },

  {
    group: "apps",
    version: "v1",
    resource: "daemonsets",
    kind: "DaemonSet",
    isNamespaced: true,
  },

  {
    group: "apps",
    version: "v1",
    resource: "statefulsets",
    kind: "StatefulSet",
    isNamespaced: true,
  },

  {
    group: "events.k8s.io",
    version: "v1",
    resource: "events",
    kind: "Event",
    isNamespaced: true,
  },

  {
    group: "batch",
    version: "v1",
    resource: "jobs",
    kind: "Job",
    isNamespaced: true,
  },

  {
    group: "batch",
    version: "v1",
    resource: "cronjobs",
    kind: "CronJob",
    isNamespaced: true,
  },

  {
    group: "networking.k8s.io",
    version: "v1",
    resource: "ingresses",
    kind: "Ingress",
    isNamespaced: true,
  },

  {
    group: "networking.k8s.io",
    version: "v1",
    resource: "ingressclasses",
    kind: "Ingressclass",
    isNamespaced: false,
  },

  {
    group: "networking.k8s.io",
    version: "v1",
    resource: "networkpolicies",
    kind: "NetworkPolicy",
    isNamespaced: false,
  },

  {
    group: "rbac.authorization.k8s.io",
    version: "v1",
    resource: "roles",
    kind: "Role",
    isNamespaced: true,
  },

  {
    group: "rbac.authorization.k8s.io",
    version: "v1",
    resource: "rolebindings",
    kind: "RoleBinding",
    isNamespaced: true,
  },

  {
    group: "rbac.authorization.k8s.io",
    version: "v1",
    resource: "clusterroles",
    kind: "ClusterRole",
    isNamespaced: false,
  },

  {
    group: "rbac.authorization.k8s.io",
    version: "v1",
    resource: "clusterrolebindings",
    kind: "ClusterRoleBinding",
    isNamespaced: false,
  },

  {
    group: "storage.k8s.io",
    version: "v1",
    resource: "storageclasses",
    kind: "StorageClass",
    isNamespaced: false,
  },
  {
    group: "apiextensions.k8s.io",
    version: "v1",
    resource: "customresourcedefinitions",
    kind: "CustomResourceDefinition",
    isNamespaced: false,
  },
];

export const NamespaceV1GVRK: KubeGVRK = {
  group: "",
  version: "v1",
  resource: "namespaces",
  kind: "Namespace",
} as const;

export const NodeV1GVRK: KubeGVRK = {
  group: "",
  version: "v1",
  resource: "nodes",
  kind: "Node",
} as const;

export const NodeMetricsV1GV: KubeGVRK = {
  group: "metrics.k8s.io",
  version: "v1beta1",
  resource: "nodes",
  kind: "NodeMetrics",
} as const;

export const PersistentVolumeV1GVRK: KubeGVRK = {
  group: "",
  version: "v1",
  resource: "persistentvolumes",
  kind: "PersistentVolume",
} as const;

export const PersistentVolumeClaimV1GVRK: KubeGVRK = {
  group: "",
  version: "v1",
  resource: "persistentvolumeclaims",
  kind: "PersistentVolumeClaim",
} as const;

export const ReplicationControllerV1GVRK: KubeGVRK = {
  group: "",
  version: "v1",
  resource: "replicationcontrollers",
  kind: "ReplicationController",
} as const;

export const PodV1GVRK: KubeGVRK = {
  group: "",
  version: "v1",
  resource: "pods",
  kind: "Pod",
} as const;

export const PodMetricsV1GV: KubeGVRK = {
  group: "metrics.k8s.io",
  version: "v1beta1",
  resource: "pods",
  kind: "PodMetrics",
} as const;

export const ServiceV1GVRK: KubeGVRK = {
  group: "",
  version: "v1",
  resource: "services",
  kind: "Service",
} as const;

export const SecretV1GVRK: KubeGVRK = {
  group: "",
  version: "v1",
  resource: "secrets",
  kind: "Secret",
} as const;

export const ServiceAccountV1GVRK: KubeGVRK = {
  group: "",
  version: "v1",
  resource: "serviceaccounts",
  kind: "ServiceAccount",
} as const;

export const ConfigMapV1GVRK: KubeGVRK = {
  group: "",
  version: "v1",
  resource: "configmaps",
  kind: "ConfigMap",
} as const;

export const DeploymentV1GVRK: KubeGVRK = {
  group: "apps",
  version: "v1",
  resource: "deployments",
  kind: "Deployment",
} as const;

export const DaemonSetV1GVRK: KubeGVRK = {
  group: "apps",
  version: "v1",
  resource: "daemonsets",
  kind: "DaemonSet",
} as const;

export const StatefulSetV1GVRK: KubeGVRK = {
  group: "apps",
  version: "v1",
  resource: "statefulsets",
  kind: "StatefulSet",
} as const;

export const EventV1GVRK: KubeGVRK = {
  group: "events.k8s.io",
  version: "v1",
  resource: "events",
  kind: "Event",
} as const;

export const JobV1GVRK: KubeGVRK = {
  group: "batch",
  version: "v1",
  resource: "jobs",
  kind: "Job",
} as const;

export const CronJobV1GVRK: KubeGVRK = {
  group: "batch",
  version: "v1",
  resource: "cronjobs",
  kind: "CronJob",
} as const;

export const IngressV1GVRK: KubeGVRK = {
  group: "networking.k8s.io",
  version: "v1",
  resource: "ingresses",
  kind: "Ingress",
} as const;

export const IngressClassV1GVRK: KubeGVRK = {
  group: "networking.k8s.io",
  version: "v1",
  resource: "ingressclasses",
  kind: "Ingressclass",
} as const;

export const NetworkPolicyV1GVRK: KubeGVRK = {
  group: "networking.k8s.io",
  version: "v1",
  resource: "networkpolicies",
  kind: "NetworkPolicy",
} as const;

export const RoleV1GVRK: KubeGVRK = {
  group: "rbac.authorization.k8s.io",
  version: "v1",
  resource: "roles",
  kind: "Role",
} as const;

export const RoleBindingV1GVRK: KubeGVRK = {
  group: "rbac.authorization.k8s.io",
  version: "v1",
  resource: "rolebindings",
  kind: "RoleBinding",
} as const;

export const ClusterRoleV1GVRK: KubeGVRK = {
  group: "rbac.authorization.k8s.io",
  version: "v1",
  resource: "clusterroles",
  kind: "ClusterRole",
} as const;

export const ClusterRoleBindingV1GVRK: KubeGVRK = {
  group: "rbac.authorization.k8s.io",
  version: "v1",
  resource: "clusterrolebindings",
  kind: "ClusterRoleBinding",
} as const;

export const StorageClassV1GVRK: KubeGVRK = {
  group: "storage.k8s.io",
  version: "v1",
  resource: "storageclasses",
  kind: "StorageClass",
} as const;

export const CustomResourceDefinitionV1GVRK: KubeGVRK = {
  group: "apiextensions.k8s.io",
  version: "v1",
  resource: "customresourcedefinitions",
  kind: "CustomResourceDefinition",
} as const;

// export const KubeResourceTypeList: Array<{ name: string; gvrk: KubeGVRK }> = [
//   { name: "ClusterRole", gvrk: ClusterRoleV1GVRK },
//   { name: "ClusterRoleBinding", gvrk: ClusterRoleBindingV1GVRK },
//   { name: "ConfigMap", gvrk: ConfigMapV1GVRK },
//   { name: "CronJob", gvrk: CronJobV1GVRK },
//   { name: "DaemonSet", gvrk: DaemonSetV1GVRK },
//   { name: "Deployment", gvrk: DeploymentV1GVRK },
//   { name: "Event", gvrk: EventV1GVRK },
//   { name: "Ingress", gvrk: IngressV1GVRK },
//   { name: "IngressClass", gvrk: IngressClassV1GVRK },
//   { name: "Job", gvrk: JobV1GVRK },
//   { name: "Namespace", gvrk: NamespaceV1GVRK },
//   { name: "NetworkPolicy", gvrk: NetworkPolicyV1GVRK },
//   { name: "Node", gvrk: NodeV1GVRK },
//   { name: "PersistentVolume", gvrk: PersistentVolumeV1GVRK },
//   { name: "PersistentVolumeClaim", gvrk: PersistentVolumeClaimV1GVRK },
//   { name: "Pod", gvrk: PodV1GVRK },
//   { name: "ReplicationController", gvrk: ReplicationControllerV1GVRK },
//   { name: "Role", gvrk: RoleV1GVRK },
//   { name: "RoleBinding", gvrk: RoleBindingV1GVRK },
//   { name: "Secret", gvrk: SecretV1GVRK },
//   { name: "Service", gvrk: ServiceV1GVRK },
//   { name: "ServiceAccount", gvrk: ServiceAccountV1GVRK },
//   { name: "StatefulSet", gvrk: StatefulSetV1GVRK },
//   { name: "StorageClass", gvrk: StorageClassV1GVRK },
// ];

// "ClusterRole",
//   "ClusterRoleBinding",
//   "ConfigMap",
//   "CronJob",
//   "DaemonSet",
//   "Deployment",
//   "Event",
//   "Ingress",
//   "IngressClass",
//   "Job",
//   "Namespace",
//   "NetworkPolicy",
//   "Node",
//   "PersistentVolume",
//   "PersistentVolumeClaim",
//   "Pod",
//   "ReplicationController",
//   "Role",
//   "RoleBinding",
//   "Secret",
//   "Service",
//   "ServiceAccount",
//   "StatefulSet",
//   "StorageClass",
