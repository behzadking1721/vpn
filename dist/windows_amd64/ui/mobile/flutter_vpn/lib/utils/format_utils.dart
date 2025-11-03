class FormatUtils {
  /// فرمت‌بندی بایت‌ها به فرمت قابل خواندن
  static String formatBytes(int bytes) {
    const suffixes = ['B', 'KB', 'MB', 'GB', 'TB'];
    var i = 0;
    var value = bytes.toDouble();
    
    while (value >= 1024 && i < suffixes.length - 1) {
      value /= 1024;
      i++;
    }
    
    return '${value.toStringAsFixed(2)} ${suffixes[i]}';
  }
  
  /// فرمت‌بندی مدت زمان به فرمت قابل خواندن
  static String formatDuration(Duration duration) {
    final hours = duration.inHours;
    final minutes = duration.inMinutes.remainder(60);
    final seconds = duration.inSeconds.remainder(60);
    
    if (hours > 0) {
      return '${hours.toString().padLeft(2, '0')}:${minutes.toString().padLeft(2, '0')}:${seconds.toString().padLeft(2, '0')}';
    } else {
      return '${minutes.toString().padLeft(2, '0')}:${seconds.toString().padLeft(2, '0')}';
    }
  }
  
  /// فرمت‌بندی پینگ به فرمت قابل خواندن
  static String formatPing(int ping) {
    if (ping <= 0) return 'N/A';
    if (ping < 50) return '$ping ms (Excellent)';
    if (ping < 100) return '$ping ms (Good)';
    if (ping < 200) return '$ping ms (Fair)';
    return '$ping ms (Poor)';
  }
}