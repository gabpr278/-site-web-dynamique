# -site-web-dynamique
Description du projet

Ce projet a été réalisé dans le cadre du TP de programmation Web avec Golang, dont l’objectif est de rendre un site e-commerce dynamique en utilisant le langage Go, le package net/http et le système de templates HTML.

Le site permet :
d’afficher une liste de produits avec leurs informations (Challenge 01) ;
de consulter le détail d’un produit sélectionné (Challenge 02) ;
d’ajouter un nouveau produit via un formulaire (Challenge 03).       



Challenge 01 — Liste des articles
Affiche tous les produits depuis une variable globale côté serveur.
Chaque article affiche : une image, un nom, une description, un prix et une éventuelle réduction.
Un bouton “Voir le produit” permet d’accéder à la page de détails.

Challenge 02 — Détail d’un article
Affiche toutes les informations d’un produit sélectionné.
Gestion des erreurs si l’article n’existe pas.
Affiche les réductions et le stock disponible.

Challenge 03 — Ajouter un produit
Formulaire d’ajout de produit avec validation des champs (nom, description, prix, etc.).
Ajout automatique à la variable globale articles.
Redirection vers la page du produit nouvellement ajouté.