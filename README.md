# fhir-api


# Registro de Decisão de Arquitetura: Serviço API

**Status:** Proposto  
**Data:** 18/08/2025  
**Autor:** Patrick Barreto


**Versão:** 1.0

---

## 1. Contexto

Necessidade de uma API robusta para:
- Gerenciar recursos FHIR (Pacientes, Profissionais, Encontros)
- Fornecer endpoints RESTful compatíveis com padrão FHIR
- Integrar com MongoDB
- Operar em ambiente cloud-native
- Manter alta disponibilidade e escalabilidade
- Observalibilidade

---

## 2. Decisão Arquitetural

### Tecnologias Principais
| Componente       | Tecnologia Escolhida      |
|------------------|---------------------------|
| Framework        | Gin                       |
| Banco de Dados   | MongoDB                   |
| Documentação     | Swagger UI                |
| Monitoramento    | Prometheus metrics        |
| Containerização  | Docker                    |
| Token            | JWT                       |


---

## 3. Consequências

### ✅ Benefícios
- **Performance:** Alta capacidade de requisições/seg
- **Observabilidade:** Métricas, logs e healthchecks
- **Segurança:** JWT + Middlewares de proteção
- **Escalabilidade:** Multi-tenant

---

## 7. Links Relevantes
- [Repositório GitHub](https://github.com/barretodsp/fhir-api)
- [FHIR R4 Specification](http://hl7.org/fhir/R4)

---

**Como usar este template:**
1. Clone o repositório
2. Execute `make docs` para gerar documentação
3. Acesse `/swagger` para API docs
4. Monitoramento disponível em `/metrics`

**Recomendação:** Manter este arquivo em `/docs/adr/001-fhir-api-architecture.md`