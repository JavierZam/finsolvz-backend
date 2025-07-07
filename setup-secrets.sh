#!/bin/bash

# Finsolvz Backend - GCP Secrets Setup
set -e

echo "🔐 Setting up GCP Secrets for Finsolvz Backend"
echo "=============================================="

# Disable bash history expansion untuk menghindari masalah dengan !
set +H

# MongoDB URI
echo "📊 Creating MONGO_URI secret..."
gcloud secrets create MONGO_URI --data-file=<(echo "mongodb+srv://tasyadviz:!xCiCF5ZyaN!E9K@adviz.sqhy6.mongodb.net/Finsolvz?retryWrites=true&w=majority")

# JWT Secret (sudah berhasil sebelumnya)
echo "🔑 JWT_SECRET already created ✅"

# Email secrets (optional untuk testing)
echo "📧 Creating email secrets..."
read -p "Enter your Gmail address (or press Enter to skip): " GMAIL_EMAIL
if [ ! -z "$GMAIL_EMAIL" ]; then
    gcloud secrets create NODEMAILER_EMAIL --data-file=<(echo "$GMAIL_EMAIL")
    
    read -s -p "Enter your Gmail App Password (or press Enter to skip): " GMAIL_PASS
    echo ""
    if [ ! -z "$GMAIL_PASS" ]; then
        gcloud secrets create NODEMAILER_PASS --data-file=<(echo "$GMAIL_PASS")
    fi
fi

# Enable bash history expansion kembali
set -H

echo ""
echo "✅ Secrets setup complete!"
echo "🔍 Verify secrets:"
gcloud secrets list

echo ""
echo "🔐 Secret versions:"
gcloud secrets versions list JWT_SECRET
gcloud secrets versions list MONGO_URI

echo ""
echo "🚀 Ready for deployment!"