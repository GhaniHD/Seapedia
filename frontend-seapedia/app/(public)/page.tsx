import { redirect } from "next/navigation";

// SEAPEDIA tidak punya halaman landing terpisah dari katalog — begitu dibuka,
// pengunjung langsung diarahkan ke /products (hero, kategori, flash sale, dan
// review aplikasi semuanya ada di sana). Ini menggantikan halaman default
// bawaan create-next-app yang sebelumnya tanpa sengaja menang di path "/".
export default function RootPage() {
  redirect("/products");
}
