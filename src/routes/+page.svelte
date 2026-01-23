<script lang="ts">
</script>

<body class="bg-gray-100 min-h-screen flex items-center justify-center p-4">
  <div
    class="border w-full max-w-sm overflow-hidden shadow-2xl bg-white rounded-3xl p-8 transition-all duration-500 hover:shadow-[0_20px_50px_rgba(0,0,0,0.2)]"
  >
    <div class="mb-8">
      <div class="flex justify-center mb-2">
        <span
          class="bg-blue-100 text-[#0065F8] text-[10px] font-bold px-2 py-1 rounded-full uppercase tracking-widest"
          >Official Portal</span
        >
      </div>
      <h1
        class="font-bold text-2xl text-center text-gray-800 font-[Poppins] tracking-tight"
      >
        Besq User Login
      </h1>
      <div
        class="text-[#0065F8] text-7xl text-center flex items-center justify-center mt-6 drop-shadow-md"
      >
        <i class="fas fa-circle-user text-blue-600"></i>
      </div>
    </div>

    <div class="space-y-5">
      <div>
        <label
          class="text-[11px] font-bold text-gray-400 uppercase tracking-widest ml-1"
          >Username</label
        >
        <div class="relative mt-1">
          <input
            type="text"
            placeholder="Masukkan username"
            class="border border-gray-100 bg-gray-50 rounded-xl w-full py-3 px-4 font-[inter] outline-none focus:bg-white focus:border-[#0065F8] focus:ring-4 focus:ring-blue-100 transition-all placeholder:text-gray-300"
          />
        </div>
      </div>

      <div>
        <label
          class="text-[11px] font-bold text-gray-400 uppercase tracking-widest ml-1"
          >Password</label
        >
        <div class="relative mt-1">
          <input
            type="password"
            placeholder="••••••••"
            class="border border-gray-100 bg-gray-50 rounded-xl w-full py-3 px-4 font-[inter] outline-none focus:bg-white focus:border-[#0065F8] focus:ring-4 focus:ring-blue-100 transition-all placeholder:text-gray-300"
          />
        </div>
      </div>
    </div>

    <button
      class="bg-[#0065F8] shadow-lg shadow-blue-300/50 transition-all duration-300 hover:bg-[#004bbd] hover:shadow-blue-400/50 hover:-translate-y-1 text-white font-bold py-3.5 px-4 rounded-xl w-full mt-8 cursor-pointer active:scale-95 flex justify-center items-center gap-2"
    >
      <span>Login</span>
    </button>

    <div class="mt-8 pt-6 border-t border-gray-100">
      <p class="text-center text-gray-400 text-[11px] leading-relaxed italic">
        Ada masalah akun? Sampaikan kendala kepada <br />
        <span class="text-gray-600 font-semibold not-italic"
          >Administrator Sistem BESQ</span
        >
      </p>
    </div>
  </div>
</body>
