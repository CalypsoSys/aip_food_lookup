import 'package:flutter/material.dart';
import 'package:google_mobile_ads/google_mobile_ads.dart';

import '../app/config.dart';

class SearchAdBanner extends StatefulWidget {
  const SearchAdBanner({
    super.key,
    this.adsEnabled = AppConfig.adsEnabled,
    this.adUnitId = AppConfig.adMobBannerAdUnitId,
  });

  final bool adsEnabled;
  final String adUnitId;

  @override
  State<SearchAdBanner> createState() => _SearchAdBannerState();
}

class _SearchAdBannerState extends State<SearchAdBanner> {
  BannerAd? _ad;
  bool _isLoaded = false;

  @override
  void initState() {
    super.initState();
    _loadAd();
  }

  @override
  void didUpdateWidget(SearchAdBanner oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (oldWidget.adsEnabled != widget.adsEnabled ||
        oldWidget.adUnitId != widget.adUnitId) {
      _disposeAd();
      _loadAd();
    }
  }

  @override
  void dispose() {
    _disposeAd();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final ad = _ad;
    if (!widget.adsEnabled || !_isLoaded || ad == null) {
      return const SizedBox.shrink();
    }

    return SizedBox(
      width: ad.size.width.toDouble(),
      height: ad.size.height.toDouble(),
      child: AdWidget(ad: ad),
    );
  }

  void _loadAd() {
    final adUnitId = widget.adUnitId.trim();
    if (!widget.adsEnabled || adUnitId.isEmpty) {
      return;
    }

    _ad = BannerAd(
      adUnitId: adUnitId,
      request: const AdRequest(),
      size: AdSize.banner,
      listener: BannerAdListener(
        onAdLoaded: (ad) {
          if (!mounted) {
            ad.dispose();
            return;
          }
          setState(() {
            _isLoaded = true;
          });
        },
        onAdFailedToLoad: (ad, error) {
          ad.dispose();
          if (!mounted) {
            return;
          }
          setState(() {
            _ad = null;
            _isLoaded = false;
          });
        },
      ),
    )..load();
  }

  void _disposeAd() {
    _ad?.dispose();
    _ad = null;
    _isLoaded = false;
  }
}
