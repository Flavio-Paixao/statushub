# ⛽ Painel de Combustíveis no Brasil

[![Streamlit App](https://img.shields.io/badge/Streamlit-Live-FF6600?style=flat-square&logo=streamlit&logoColor=white)](https://combustiveis-brasil-anp-gu8tg5fercq66pecnwfyp6.streamlit.app/)
[![Python](https://img.shields.io/badge/Python-3.11+-3776AB?style=flat-square&logo=python&logoColor=white)](https://python.org)
[![Pandas](https://img.shields.io/badge/Pandas-2.x-150458?style=flat-square&logo=pandas&logoColor=white)](https://pandas.pydata.org)
[![Plotly](https://img.shields.io/badge/Plotly-Interactive-3F4F75?style=flat-square&logo=plotly&logoColor=white)](https://plotly.com)
[![ANP](https://img.shields.io/badge/Fonte-ANP%20Dados%20Abertos-22C55E?style=flat-square)](https://www.gov.br/anp)

> Dashboard interativo para análise exploratória dos preços de revenda de combustíveis no Brasil, com dados oficiais da Agência Nacional do Petróleo (ANP) — série histórica semanal por estado de 2012 a 2026.

---

## 🚀 Demo ao vivo

**[→ Acessar o dashboard](https://combustiveis-brasil-anp-gu8tg5fercq66pecnwfyp6.streamlit.app/)**

---

## 📊 O que o dashboard analisa

| Seção | Descrição |
|---|---|
| **Visão Geral (KPIs)** | Preço médio, mínimo, máximo e total de registros no período filtrado |
| **Preço por Estado** | Ranking de todos os estados ordenado do mais barato ao mais caro |
| **Série Histórica** | Evolução do preço médio nacional de 2012 a 2026 |
| **Regra dos 70%** | Comparativo econômico Etanol vs Gasolina por estado |
| **Comparativo Regional** | Evolução de preços por região ao longo dos anos |
| **Top 5 Ranking** | Os 5 estados mais baratos e os 5 mais caros |

---

## 💡 Por que esse projeto?

Trabalhei como motorista de aplicativo e sempre precisei tomar decisões sobre combustível: abastecer com etanol ou gasolina? Em qual estado o preço compensa mais? Esse dashboard responde essas perguntas com dados reais da ANP.

A **Regra dos 70%** é o insight central: se o preço do etanol for menor que 70% do preço da gasolina, o etanol é mais econômico por quilômetro rodado. O dashboard identifica automaticamente quais estados estão abaixo ou acima desse limiar.

---

## 🛠️ Stack

```
Python 3.11+
├── pandas          → limpeza e análise dos dados
├── plotly          → gráficos interativos
├── streamlit       → dashboard web
└── requests        → carregamento dos dados via URL
```

**Dados:** ANP — Série Histórica Semanal por Estado (2012–2026)  
**Deploy:** Streamlit Cloud (CI/CD automático via GitHub)

---

## 📁 Estrutura do projeto

```
combustiveis-brasil-anp/
├── app.py                          # Dashboard Streamlit
├── dados_anp.csv                   # Dataset limpo (114k registros)
├── requirements.txt                # Dependências Python
├── .streamlit/
│   └── config.toml                 # Tema dark forçado
└── README.md
```

---

## ⚙️ Como rodar localmente

```bash
# 1. Clone o repositório
git clone https://github.com/Flavio-Paixao/combustiveis-brasil-anp.git
cd combustiveis-brasil-anp

# 2. Instale as dependências
pip install -r requirements.txt

# 3. Rode o dashboard
streamlit run app.py
```

---

## 📦 Fonte dos dados

Os dados são provenientes da **ANP — Agência Nacional do Petróleo, Gás Natural e Biocombustíveis**, disponíveis publicamente em:

🔗 [Série Histórica do Levantamento de Preços](https://www.gov.br/anp/pt-br/assuntos/precos-e-defesa-da-concorrencia/precos/precos-revenda-e-de-distribuicao-combustiveis/serie-historica-do-levantamento-de-precos)

**Combustíveis disponíveis no dataset:**
- Gasolina Comum e Aditivada
- Etanol Hidratado
- Óleo Diesel e Diesel S10
- GLP (Gás de Cozinha)
- GNV

---

## 🗺️ Próximos projetos da série

Este é o **Projeto 1** de uma série temática sobre mobilidade e energia no Brasil:

| # | Projeto | Stack | Status |
|---|---|---|---|
| 1 | **Painel de Combustíveis ANP** | Python · Pandas · Plotly · Streamlit | ✅ Live |
| 2 | Dashboard de Mobilidade Urbana SP | Python · Streamlit · Mapbox | 🔧 Em breve |
| 3 | ML: Previsão de Preço de Combustível | Scikit-learn · Random Forest | 🔧 Em breve |

---

## 👤 Autor

**Flávio Paixão** — Backend Developer · AWS Cloud · Data Analytics

[![LinkedIn](https://img.shields.io/badge/LinkedIn-flaviopx-0A66C2?style=flat-square&logo=linkedin)](https://linkedin.com/in/flaviopx)
[![GitHub](https://img.shields.io/badge/GitHub-Flavio--Paixao-181717?style=flat-square&logo=github)](https://github.com/Flavio-Paixao)
[![Instagram](https://img.shields.io/badge/Instagram-flaviopaixao.dev-E4405F?style=flat-square&logo=instagram)](https://instagram.com/flaviopaixao.dev)
