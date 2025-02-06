#!/bin/bash

# Vérifier si des fichiers SVG existent
if ls *.svg &>/dev/null; then
    # Trier les fichiers en ordre décroissant pour éviter les conflits (ex: 18.svg avant 17.svg)
    for file in $(ls *.svg | sort -Vr); do
        # Extraire le numéro du fichier
        num=$(basename "$file" .svg)
        
        # Vérifier si le numéro est supérieur à 16
        if [[ $num -gt 16 ]]; then
            # Calculer le nouveau numéro
            new_num=$((num + 1))
            
            # Renommer le fichier
            mv "$file" "${new_num}.svg"
        fi
    done
    echo "Renommage terminé !"
else
    echo "Aucun fichier SVG trouvé."
fi
