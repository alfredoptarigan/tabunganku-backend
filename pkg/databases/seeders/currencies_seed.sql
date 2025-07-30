-- Seed data untuk currencies
INSERT INTO currencies (country_name, currency_name, country_flag, currency_symbol, currency_code) VALUES
('Indonesia', 'Indonesian Rupiah', '🇮🇩', 'Rp', 'IDR'),
('United States', 'US Dollar', '🇺🇸', '$', 'USD'),
('European Union', 'Euro', '🇪🇺', '€', 'EUR'),
('United Kingdom', 'British Pound', '🇬🇧', '£', 'GBP'),
('Japan', 'Japanese Yen', '🇯🇵', '¥', 'JPY'),
('Singapore', 'Singapore Dollar', '🇸🇬', 'S$', 'SGD'),
('Malaysia', 'Malaysian Ringgit', '🇲🇾', 'RM', 'MYR'),
('Thailand', 'Thai Baht', '🇹🇭', '฿', 'THB'),
('South Korea', 'South Korean Won', '🇰🇷', '₩', 'KRW'),
('China', 'Chinese Yuan', '🇨🇳', '¥', 'CNY'),
('Australia', 'Australian Dollar', '🇦🇺', 'A$', 'AUD'),
('Canada', 'Canadian Dollar', '🇨🇦', 'C$', 'CAD'),
('Switzerland', 'Swiss Franc', '🇨🇭', 'CHF', 'CHF'),
('India', 'Indian Rupee', '🇮🇳', '₹', 'INR'),
('Brazil', 'Brazilian Real', '🇧🇷', 'R$', 'BRL')
ON CONFLICT (currency_code) DO NOTHING;