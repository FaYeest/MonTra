/// Primary navigation destinations defined in `docs/design.md` Section 4.1.
///
/// Persistent bottom navigation destinations:
/// 1. Beranda — financial overview and recent transactions
/// 2. Transaksi — searchable transaction history
/// 3. Tambah — central, visually emphasized action to record a transaction
/// 4. Anggaran — overall and category budgets
/// 5. Insight — trends, forecasts, and explanations
///
/// Profile (`Profil`) is accessed via the avatar in the top app bar and is not a sixth bottom item.
enum MontraPrimaryTab {
  beranda(label: 'Beranda', semanticLabel: 'Tab Beranda, ringkasan keuangan'),
  transaksi(
    label: 'Transaksi',
    semanticLabel: 'Tab Transaksi, riwayat transaksi',
  ),
  anggaran(
    label: 'Anggaran',
    semanticLabel: 'Tab Anggaran, batas pengeluaran',
  ),
  insight(
    label: 'Insight',
    semanticLabel: 'Tab Insight, analisis dan perkiraan',
  );

  const MontraPrimaryTab({
    required this.label,
    required this.semanticLabel,
  });

  final String label;
  final String semanticLabel;
}
