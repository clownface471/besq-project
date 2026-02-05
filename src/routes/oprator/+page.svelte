	<script lang="ts">
		// @ts-nocheck
		import { goto } from '$app/navigation';
		import { auth, logout } from '$lib/stores/auth';
		import Swal from 'sweetalert2';

		// Helper untuk cek izin (case insensitive biar aman)
	const canAccess = (roleRequired: string) => {
		const userRole = $auth.user?.role || '';
		if (userRole === 'ADMIN') return true; // Admin bisa semua
		return userRole.includes(roleRequired);
	};

	// Function navigasi dengan proteksi
	const handleNavigate = (path: string, requiredRole: string) => {
		if (canAccess(requiredRole)) {
			goto(path);
		} else {
			Swal.fire({
				toast: true,
				icon: 'error',
				title: 'Akses Ditolak',
				text: `Anda bukan operator ${requiredRole.toLowerCase()}`,
				position: 'top-end',
				timer: 2000,
				showConfirmButton: false
			});
		}
	};

		// --- PROTECT ROUTE ---
		// Jika user belum login, lempar kembali ke halaman login
		$effect(() => {
			if (!$auth.isLoggedIn) {
				goto('/');
			}
		});

		// --- DYNAMIC DATA ---
		// Mengambil data real dari user yang sedang login
let employee = $derived({
		name: $auth.user?.name || 'User',
		position: $auth.user?.role || 'Staff',
		department: 'Production',
		nikId: $auth.user?.username || '-',
		phone: '-',
		shift: 'Shift 1',
		avatar: `https://ui-avatars.com/api/?name=${$auth.user?.name}&background=random`,
		date: new Date().toLocaleDateString('id-ID', { weekday: 'long', day: 'numeric', month: 'long' })
	});

		const handleLogout = async () => {
			const result = await Swal.fire({
				title: 'Logout?',
				text: 'Anda harus login ulang untuk masuk kembali.',
				icon: 'warning',
				showCancelButton: true,
				confirmButtonText: 'Ya, Keluar',
				confirmButtonColor: '#d33',
				cancelButtonText: 'Batal'
			});

			if (result.isConfirmed) {
				logout();
			}
		};
	</script>

	<div class="container">
		<div class="absolute top-4 right-4 z-10">
			<button on:click={handleLogout} class="bg-white/80 hover:bg-white text-red-500 px-4 py-2 rounded-xl shadow-sm font-bold text-xs transition-all">
				LOGOUT
			</button>
		</div>

		<div class="profile-section">
			<div class="profile-left">
				<img src={employee.avatar} alt={employee.name} class="avatar" />
			</div>

			<div class="profile-center">
				<h1>{employee.name}</h1>
				<div class="badges">
					<span class="badge badge-position">{employee.position}</span>
					<span class="badge badge-department">{employee.department}</span>
				</div>

				<div class="profile-details">
					<div class="detail-item">
						<span class="detail-icon">📋</span>
						<div>
							<p class="detail-label">NIK ID</p>
							<p class="detail-value uppercase">{employee.nikId}</p>
						</div>
					</div>

					<div class="detail-item">
						<span class="detail-icon">☎️</span>
						<div>
							<p class="detail-label">TELEPON</p>
							<p class="detail-value">{employee.phone}</p>
						</div>
					</div>

					<div class="detail-item">
						<span class="detail-icon">🕐</span>
						<div>
							<p class="detail-label">SHIFT</p>
							<p class="detail-value">{employee.shift}</p>
						</div>
					</div>
				</div>
			</div>

			<div class="profile-right">
				<p class="date-label">Hari ini</p>
				<p class="date-value">{employee.date}</p>
			</div>
		</div>

		<div class="main-content">
			<div class="boxes-wrapper">
				<div class="boxes-grid">
					<div class="box-wrapper">
						<div 
						class="box cutting"
						class:opacity-50={!canAccess('CUTTING')}
						class:cursor-not-allowed={!canAccess('CUTTING')}
						on:click={() => handleNavigate('/cutting', 'CUTTING')}
						role="button" 
						tabindex="0"
					>
							<svg class="box-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
							<circle cx="6" cy="6" r="3"></circle>
							<circle cx="6" cy="18" r="3"></circle>
							<line x1="20" y1="4" x2="9" y2="15"></line>
							<line x1="20" y1="20" x2="9" y2="9"></line>
						</svg>
					</div>
					<p class="box-label">Cutting</p>
				</div>
					
					<div class="connector-line"></div>

					<div class="box-wrapper">
<div 
						class="box pressing" 
						class:opacity-50={!canAccess('PRESSING')}
						class:cursor-not-allowed={!canAccess('PRESSING')}
						on:click={() => handleNavigate('/pressing', 'PRESSING')}
						role="button" 
						tabindex="0"
					>
							<svg class="box-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
								<path d="M12 3v6m0 6v6"></path>
								<path d="M8 10h8"></path>
								<path d="M8 20h8"></path>
								<path d="M6 14h12"></path>
							</svg>
						</div>
						<p class="box-label">Pressing</p>
					</div>
				</div>
			</div>
		</div>
	</div>

	<style>
		/* ... (Style tetap sama seperti sebelumnya, tidak perlu diubah) ... */
		.container {
			min-height: 100vh;
			width: 100%;
			background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
			padding: 2rem;
			display: flex;
			flex-direction: column;
			position: relative; /* Penting untuk tombol logout absolute */
		}
		/* Sisanya copy-paste style dari file lama Anda */
		/* Profile Section */
		.profile-section {
			display: flex;
			align-items: center;
			gap: 2rem;
			background: white;
			padding: 1.5rem 2rem;
			border-radius: 16px;
			box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
			margin-bottom: 3rem;
			max-width: 1400px;
			margin-left: auto;
			margin-right: auto;
		}

		.profile-left {
			flex-shrink: 0;
		}

		.avatar {
			width: 100px;
			height: 100px;
			border-radius: 16px;
			object-fit: cover;
			box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
		}

		.profile-center {
			flex: 1;
		}

		.profile-center h1 {
			margin: 0 0 0.8rem 0;
			font-size: 1.8rem;
			color: #1a202c;
			font-weight: 700;
		}

		.badges {
			display: flex;
			gap: 0.8rem;
			margin-bottom: 1.2rem;
			flex-wrap: wrap;
		}

		.badge {
			padding: 0.4rem 0.9rem;
			border-radius: 16px;
			font-size: 0.7rem;
			font-weight: 600;
			letter-spacing: 0.5px;
			display: inline-block;
		}

		.badge-position {
			background: #e0e7ff;
			color: #5a67d8;
		}

		.badge-department {
			background: #f3f4f6;
			color: #6b7280;
		}

		.profile-details {
			display: flex;
			gap: 2rem;
			flex-wrap: wrap;
		}

		.detail-item {
			display: flex;
			align-items: center;
			gap: 0.8rem;
		}

		.detail-icon {
			font-size: 1.3rem;
			flex-shrink: 0;
		}

		.detail-label {
			margin: 0;
			font-size: 0.65rem;
			color: #9ca3af;
			font-weight: 600;
			letter-spacing: 0.5px;
			text-transform: uppercase;
		}

		.detail-value {
			margin: 0.2rem 0 0 0;
			font-size: 0.85rem;
			color: #1f2937;
			font-weight: 500;
		}

		.profile-right {
			flex-shrink: 0;
			text-align: right;
			padding-left: 2rem;
			border-left: 1px solid #e5e7eb;
		}

		.date-label {
			margin: 0;
			font-size: 0.75rem;
			color: #9ca3af;
			font-weight: 500;
		}

		.date-value {
			margin: 0.4rem 0 0 0;
			font-size: 1.1rem;
			color: #1a202c;
			font-weight: 700;
		}

		/* Main Content */
		.main-content {
			max-width: 1400px;
			margin: 0 auto;
		}

		.boxes-wrapper {
			display: flex;
			justify-content: center;
			align-items: center;
		}

		.boxes-grid {
			display: grid;
			grid-template-columns: auto auto auto;
			gap: 0;
			align-items: center;
			padding: 0 1rem;
			position: relative;
		}

		.box-wrapper {
			display: flex;
			flex-direction: column;
			align-items: center;
			gap: 1rem;
		}

		.box {
			width: 100%;
			max-width: 300px;
			aspect-ratio: 1 / 1;
			padding: 2rem;
			border-radius: 24px;
			cursor: pointer;
			display: flex;
			justify-content: center;
			align-items: center;
			transition: all 0.4s cubic-bezier(0.23, 1, 0.320, 1);
			box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
			position: relative;
			overflow: hidden;
			background: #f3f4f6;
			border: 2px solid transparent;
		}

		.box::before {
			content: '';
			position: absolute;
			top: 0;
			left: -100%;
			width: 100%;
			height: 100%;
			background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.3), transparent);
			transition: left 0.5s ease;
		}

		.box:hover {
			transform: translateY(-8px) scale(1.02);
			box-shadow: 0 20px 40px rgba(0, 0, 0, 0.15);
		}

		.box:hover::before {
			left: 100%;
		}

		.box-icon {
			width: 80px;
			height: 80px;
			transition: transform 0.4s ease;
			color: #374151;
			position: relative;
			z-index: 1;
		}

		.box:hover .box-icon {
			transform: scale(1.15) rotate(5deg);
		}

		.box-label {
			margin: 0;
			font-size: 1rem;
			color: #1f2937;
			font-weight: 600;
			letter-spacing: 1px;
		}

		.cutting {
			background: linear-gradient(135deg, #e8eaed 0%, #f3f4f6 100%);
			border-color: #d4d7dc;
		}

		.cutting:hover {
			background: linear-gradient(135deg, #d4d7dc 0%, #e8eaed 100%);
			border-color: #bcc1c9;
		}

		.pressing {
			background: linear-gradient(135deg, #f0f1f4 0%, #f5f6f8 100%);
			border-color: #dfe2e8;
		}

		.pressing:hover {
			background: linear-gradient(135deg, #dfe2e8 0%, #e8eaed 100%);
			border-color: #cbd2db;
		}

		.connector-line {
			width: 60px;
			height: 4px;
			background: linear-gradient(90deg, #cbd5e0 0%, #a0aec0 50%, #cbd5e0 100%);
			position: relative;
			margin: 0 1rem;
		}

		.connector-line::after {
			content: '';
			position: absolute;
			right: -12px;
			top: 50%;
			transform: translateY(-50%);
			width: 0;
			height: 0;
			border-left: 12px solid #a0aec0;
			border-top: 6px solid transparent;
			border-bottom: 6px solid transparent;
		}

		/* Responsive */
		@media (max-width: 1024px) {
			.profile-section {
				flex-wrap: wrap;
				gap: 1.5rem;
			}

			.profile-right {
				border-left: none;
				border-top: 1px solid #e5e7eb;
				padding-left: 0;
				padding-top: 1rem;
				width: 100%;
				text-align: center;
			}

			.boxes-grid {
				gap: 0;
			}

			.box {
				max-width: 280px;
			}

			.box-icon {
				width: 70px;
				height: 70px;
			}

			.connector-line {
				width: 50px;
				margin: 0 0.8rem;
			}
		}

		@media (max-width: 768px) {
			.container {
				padding: 0;
				min-height: 100vh;
			}

			.profile-section {
				flex-direction: column;
				padding: 1.2rem;
				gap: 1rem;
				margin-bottom: 2rem;
				text-align: center;
			}

			.profile-left {
				flex-shrink: 0;
			}

			.avatar {
				width: 90px;
				height: 90px;
			}

			.profile-center {
				width: 100%;
			}

			.profile-center h1 {
				font-size: 1.4rem;
				margin-bottom: 0.6rem;
			}

			.badges {
				justify-content: center;
				margin-bottom: 1rem;
			}

			.badge {
				font-size: 0.65rem;
				padding: 0.3rem 0.8rem;
			}

			.profile-details {
				flex-direction: column;
				gap: 0.8rem;
				width: 100%;
			}

			.detail-item {
				justify-content: center;
			}

			.profile-right {
				border: none;
				padding: 0;
				width: auto;
			}

			.date-label {
				font-size: 0.7rem;
			}

			.date-value {
				font-size: 1rem;
			}

			.boxes-grid {
				grid-template-columns: 1fr;
				gap: 1.5rem;
				padding: 0;
				max-width: 220px;
				margin: 0 auto;
			}

			.box-wrapper {
				width: 100%;
			}

			.connector-line {
				width: 3px;
				height: 45px;
				background: linear-gradient(180deg, #cbd5e0 0%, #a0aec0 50%, #cbd5e0 100%);
				margin: 0;
			}

			.connector-line::after {
				right: auto;
				top: auto;
				bottom: -12px;
				left: 50%;
				transform: translateX(-50%);
				border-left: 6px solid transparent;
				border-right: 6px solid transparent;
				border-top: 12px solid #a0aec0;
				border-bottom: none;
			}

			.box {
				max-width: none;
				aspect-ratio: 3 / 4;
				padding: 1.2rem 0.8rem;
			}

			.box-icon {
				width: 55px;
				height: 55px;
			}

			.box-label {
				font-size: 0.85rem;
				margin-top: 0.5rem;
			}

			.box:hover .box-icon {
				transform: scale(1.1) rotate(3deg);
			}
		}

		@media (max-width: 480px) {
			.container {
				padding: 0;
				min-height: 100vh;
			}

			.profile-section {
				flex-direction: column;
				padding: 1rem;
				gap: 0.8rem;
				margin-bottom: 1.5rem;
				border-radius: 12px;
			}

			.avatar {
				width: 85px;
				height: 85px;
				border-radius: 12px;
			}

			.profile-center h1 {
				font-size: 1.3rem;
				margin-bottom: 0.5rem;
			}

			.badges {
				gap: 0.5rem;
				margin-bottom: 0.8rem;
			}

			.badge {
				font-size: 0.6rem;
				padding: 0.25rem 0.7rem;
			}

			.profile-details {
				gap: 0.6rem;
			}

			.detail-item {
				gap: 0.6rem;
				font-size: 0.85rem;
			}

			.detail-icon {
				font-size: 1.1rem;
			}

			.detail-label {
				font-size: 0.6rem;
			}

			.detail-value {
				font-size: 0.8rem;
			}

			.date-label {
				font-size: 0.65rem;
			}

			.date-value {
				font-size: 0.95rem;
			}

			.boxes-grid {
				max-width: 180px;
				gap: 1.2rem;
			}

			.connector-line {
				width: 3px;
				height: 40px;
				background: linear-gradient(180deg, #cbd5e0 0%, #a0aec0 50%, #cbd5e0 100%);
				margin: 0;
			}

			.connector-line::after {
				right: auto;
				top: auto;
				bottom: -12px;
				left: 50%;
				transform: translateX(-50%);
				border-left: 6px solid transparent;
				border-right: 6px solid transparent;
				border-top: 12px solid #a0aec0;
				border-bottom: none;
			}

			.box {
				aspect-ratio: 3 / 4;
				padding: 1rem 0.7rem;
				border-radius: 16px;
			}

			.box-icon {
				width: 50px;
				height: 50px;
			}

			.box-label {
				font-size: 0.75rem;
			}

			.box:hover .box-icon {
				transform: scale(1.05) rotate(2deg);
			}
		}

		@media (max-width: 360px) {
			.container {
				padding: 0;
				min-height: 100vh;
			}

			.profile-section {
				padding: 0.9rem;
				gap: 0.6rem;
			}

			.avatar {
				width: 75px;
				height: 75px;
			}

			.profile-center h1 {
				font-size: 1.2rem;
			}

			.boxes-grid {
				max-width: 160px;
				gap: 1rem;
			}

			.box {
				padding: 0.8rem 0.6rem;
			}
			.box-icon {
				width: 45px;
				height: 45px;
			}

			.connector-line {
				height: 35px;
			}
			.connector-line::after {
				border-left: 5px solid transparent;
				border-right: 5px solid transparent;
				border-top: 10px solid #a0aec0;
			}
		}
	</style>